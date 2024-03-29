package cmd

import (
	"bytes"
	"context"
	"fmt"
	"io"

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

func yamlReadAndSplit(reader io.Reader) ([][]byte, error) {
	fb, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	var objs [][]byte

	for _, obj := range bytes.Split(fb, []byte("---")) {
		if len(obj) != 0 {
			objs = append(objs, obj)
		}
	}

	return objs, err
}

func k3autoDeploy(ctx context.Context, name string, directory string, filesystem afero.Fs) error {
	imageRef := fmt.Sprintf("%v:%v", name, name)

	err := docker.BuildImage(ctx, directory, docker.DumpDockerfile, []string{imageRef}, filesystem)
	if err != nil {
		return err
	}

	err = docker.PushImage(ctx, imageRef, "127.0.0.1:8888")
	if err != nil {
		return err
	}

	k8sC, err := k8s.NewClient()
	if err != nil {
		return err
	}

	repo, kustomization := flux.NewOCIKustomization(name, name, name)

	err = k8sC.Create(ctx, &repo)
	if err != nil {
		return err
	}

	err = k8sC.Create(ctx, &kustomization)
	if err != nil {
		return err
	}

	logrus.Info("Default Deployments Injected")
	return nil
}
