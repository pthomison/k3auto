package k3auto

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/davecgh/go-spew/spew"
	"github.com/pthomison/k3auto/internal/containers"
	"github.com/pthomison/k3auto/internal/docker"
	"github.com/pthomison/k3auto/internal/flux"
	"github.com/pthomison/k3auto/internal/k8s"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"

	kustomizev1 "github.com/fluxcd/kustomize-controller/api/v1"
	sourcev1beta2 "github.com/fluxcd/source-controller/api/v1beta2"
)

func ensureDeployment(ctx context.Context, k8sC client.Client, name string, namespace string, repository string, image string, tag string, path string) error {
	oci := sourcev1beta2.OCIRepository{}
	err := k8sC.Get(ctx, client.ObjectKey{
		Name:      name,
		Namespace: namespace,
	}, &oci)

	if err != nil && !errors.IsNotFound(err) {
		spew.Dump(errors.IsNotFound(err), err)
		return err
	} else if err != nil {
		oci = flux.NewOCIRepoObject(name, namespace, repository, image, tag)
		err = k8sC.Create(ctx, &oci)
		if err != nil {
			return err
		}
	} else {
		oci.Spec.Reference.Tag = tag
		err = k8sC.Update(ctx, &oci)
		if err != nil {
			return err
		}
	}

	kustomization := kustomizev1.Kustomization{}
	err = k8sC.Get(ctx, client.ObjectKey{
		Name:      name,
		Namespace: namespace,
	}, &kustomization)

	if err != nil && !errors.IsNotFound(err) {
		return err
	} else if err != nil {
		kustomization = flux.NewOCIKustomizationObject(name, namespace, path)
		err = k8sC.Create(ctx, &kustomization)
		if err != nil {
			return err
		}
	}

	err = flux.WaitForKustomization(ctx, k8sC, v1.ObjectMeta{
		Name:      name,
		Namespace: namespace,
	})
	if err != nil {
		return err
	}

	return nil
}

func Deploy(ctx context.Context, name string, directory string, bootstrap string, filesystem afero.Fs) error {

	k8sC, err := k8s.NewClient()
	if err != nil {
		return err
	}

	repository := fmt.Sprintf("%v:5000", "docker-registry.docker-registry.svc.cluster.local")
	localRepository := fmt.Sprintf("%v:5000", "127.0.0.1")

	image := name
	namespace := "kube-system"

	logrus.Infof("%v Deployments Injecting", name)

	initImageRef := fmt.Sprintf("%v:%v", name, name)

	hash, err := docker.BuildImage(ctx, directory, docker.DumpDockerfile, []string{initImageRef}, filesystem)
	if err != nil {
		return err
	}

	tag := hash
	imageRef := fmt.Sprintf("%v:%v", name, hash)

	err = docker.TagImage(ctx, initImageRef, imageRef)
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

	err = ensureDeployment(ctx, k8sC, name, namespace, repository, image, tag, bootstrap)
	if err != nil {
		return err
	}

	logrus.Infof("%v Deployments Injected", name)

	return nil
}
