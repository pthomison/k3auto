package k8s

import (
	"context"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func WaitForDeployment(ctx context.Context, k8s client.Client, desiredDep v1.ObjectMeta) error {
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
