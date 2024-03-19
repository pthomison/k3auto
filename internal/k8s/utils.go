package k8s

import (
	"context"
	"os"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	apimachinerytypes "k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"k8s.io/client-go/kubernetes/scheme"

	kustomizev1 "github.com/fluxcd/kustomize-controller/api/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
)

func DeploymentReady(ctx context.Context, k8s client.Client, name string, namespace string) (bool, error) {
	dp := appsv1.Deployment{}
	err := k8s.Get(ctx, apimachinerytypes.NamespacedName{
		Name:      name,
		Namespace: namespace,
	}, &dp)
	if err != nil {
		return false, err
	}

	return (dp.Status.ReadyReplicas == *dp.Spec.Replicas), nil
}

func WaitForDeployment(ctx context.Context, k8s client.Client, desiredDep metav1.ObjectMeta) error {
	for {
		deploymentList := appsv1.DeploymentList{}
		err := k8s.List(ctx, &deploymentList, &client.ListOptions{
			Namespace: desiredDep.Namespace,
		})
		if err != nil {
			return err
		}

		for _, dep := range deploymentList.Items {
			if dep.Name == desiredDep.Name {
				if dep.Status.ReadyReplicas == *dep.Spec.Replicas {
					return nil
				}
			}
		}

		time.Sleep(1 * time.Second)
	}
}

func ParseManifestFile(fileLocation string) (runtime.Object, *schema.GroupVersionKind, error) {
	b, err := os.ReadFile(fileLocation)
	if err != nil {
		return nil, nil, err
	}

	decoder := NewDecoder()
	runtimeObject, groupVersionKind, err := decoder.Decode(b, nil, nil)
	if err != nil {
		return nil, nil, err
	}

	return runtimeObject, groupVersionKind, nil
}

func NewDecoder() runtime.Decoder {
	apiextensionsv1.AddToScheme(scheme.Scheme)
	apiextensionsv1beta1.AddToScheme(scheme.Scheme)
	kustomizev1.AddToScheme(scheme.Scheme)

	decoder := scheme.Codecs.UniversalDeserializer()

	return decoder
}
