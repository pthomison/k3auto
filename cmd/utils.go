package cmd

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/pthomison/k3auto/internal/containers"
	"github.com/pthomison/k3auto/internal/docker"
	"github.com/pthomison/k3auto/internal/flux"
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

// func lookupIpv4() (string, error) {
// 	// get list of available addresses
// 	addr, err := net.InterfaceAddrs()
// 	if err != nil {
// 		return "", err
// 	}

// 	for _, addr := range addr {
// 		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
// 			// check if IPv4 or IPv6 is not nil
// 			if ipnet.IP.To4() != nil {
// 				// print available addresses
// 				return ipnet.IP.String(), nil
// 			}
// 		}
// 	}
// 	return "", errors.New("no ipv4 network detected")
// }

func Deploy(ctx context.Context, name string, directory string, filesystem afero.Fs) error {
	imageRef := fmt.Sprintf("%v:%v", name, name)

	k8sC, err := k8s.NewClient()
	if err != nil {
		return err
	}

	tag := name

	repository := fmt.Sprintf("%v:5000", "docker-registry.docker-registry.svc.cluster.local")
	localRepository := fmt.Sprintf("%v:5000", "127.0.0.1")

	image := name
	namespace := "kube-system"

	logrus.Infof("%v Deployments Injecting", name)

	err = docker.BuildImage(ctx, directory, docker.DumpDockerfile, []string{imageRef}, filesystem)
	if err != nil {
		return err
	}

	dep := appsv1.Deployment{}
	err = k8sC.Get(ctx, client.ObjectKey{
		Name:      "docker-registry",
		Namespace: "docker-registry",
	}, &dep)
	if err != nil {
		return err
	}

	pods := corev1.PodList{}
	var selector client.MatchingLabels = dep.Spec.Selector.MatchLabels
	err = k8sC.List(ctx, &pods, selector)
	if err != nil {
		return err
	}

	closeChan, err := k8s.PortForward(ctx, pods.Items[0].Name, pods.Items[0].Namespace, 5000)
	if err != nil {
		return err
	}

	err = containers.PushImage(ctx, imageRef, localRepository)
	if err != nil {
		return err
	}

	close(closeChan)

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
