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

	return ParseManifest(b)
}

func ParseManifest(manifestData []byte) (runtime.Object, *schema.GroupVersionKind, error) {
	decoder := NewDecoder()
	runtimeObject, groupVersionKind, err := decoder.Decode(manifestData, nil, nil)
	if err != nil {
		return nil, nil, err
	}

	return runtimeObject, groupVersionKind, nil
}
