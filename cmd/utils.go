package cmd

import (
	"context"
	"errors"
	"fmt"
	"net"

	"k8s.io/apimachinery/pkg/runtime"

	k3dv1alpha5 "github.com/k3d-io/k3d/v5/pkg/config/v1alpha5"
	defaults "github.com/pthomison/k3auto/default"
	"github.com/pthomison/k3auto/internal/docker"
	"github.com/pthomison/k3auto/internal/flux"
	"github.com/pthomison/k3auto/internal/k3d"
	"github.com/pthomison/k3auto/internal/k8s"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

func checkError(err error) {
	if err != nil {
		logrus.Fatal(err)
		panic(err)
	}
}

func parseConfigFile(configPath string) (*k3dv1alpha5.SimpleConfig, error) {
	var clusterConfig *k3dv1alpha5.SimpleConfig
	var err error

	if configPath != "" {
		clusterConfig, err = k3d.ParseConfigFile(configPath, nil)
		if err != nil {
			return nil, err
		}
	} else {
		clusterConfig, err = k3d.ParseConfigFile(defaults.K3dConfigLocation, &defaults.K3dConfig)
		if err != nil {
			return nil, err
		}
	}

	return clusterConfig, nil
}

func lookupIpv4() (string, error) {
	// get list of available addresses
	addr, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, addr := range addr {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			// check if IPv4 or IPv6 is not nil
			if ipnet.IP.To4() != nil {
				// print available addresses
				return ipnet.IP.String(), nil
			}
		}
	}
	return "", errors.New("no ipv4 network detected")
}

func k3autoDeploy(ctx context.Context, name string, directory string, filesystem afero.Fs) error {
	imageRef := fmt.Sprintf("%v:%v", name, name)

	machineIP, err := lookupIpv4()
	if err != nil {
		return err
	}

	tag := name
	repository := fmt.Sprintf("%v:8888", machineIP)
	localRepository := fmt.Sprintf("%v:8888", "127.0.0.1")
	image := name
	namespace := "kube-system"

	logrus.Infof("%v Deployments Injecting", name)

	err = docker.BuildImage(ctx, directory, docker.DumpDockerfile, []string{imageRef}, filesystem)
	if err != nil {
		return err
	}

	err = docker.PushImage(ctx, imageRef, localRepository)
	if err != nil {
		return err
	}

	k8sC, err := k8s.NewClient()
	if err != nil {
		return err
	}

	repo := flux.NewOCIRepoObject(name, namespace, repository, image, tag)
	kustomization := flux.NewOCIKustomizationObject(name, namespace)

	err = k8s.CreateObjects(ctx, k8sC, []runtime.Object{&repo, &kustomization})
	if err != nil {
		return err
	}

	logrus.Infof("%v Deployments Injected", name)

	logrus.Infof("Waiting For %v Kustomization", name)
	err = flux.WaitForKustomization(ctx, k8sC, kustomization.ObjectMeta)
	if err != nil {
		return err
	}

	return nil
}
