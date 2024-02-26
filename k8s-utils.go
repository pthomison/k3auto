package main

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

// func WaitForPodsReadyInCluster(ctx context.Context, k8s client.Client) error {
// 	ready := false
// 	var err error
// 	for {
// 		ready, err = ArePodsReadyInCluster(ctx, k8s)
// 		if err != nil {
// 			return err
// 		}
// 		if ready {
// 			break
// 		}

// 		logrus.Info("Waiting on cluster")
// 		time.Sleep(1 * time.Second)
// 	}
// 	return nil
// }

// func ArePodsReadyInCluster(ctx context.Context, k8s client.Client) (bool, error) {
// 	podList := corev1.PodList{}
// 	err := k8s.List(ctx, &podList)
// 	if err != nil {
// 		return false, err
// 	}

// 	clusterReady := true
// 	unready := []string{}

// 	for _, pod := range podList.Items {
// 		podReady := IsPodReady(pod)

// 		if !podReady {
// 			clusterReady = false
// 			unready = append(unready, pod.ObjectMeta.Name)
// 		}
// 	}

// 	if !clusterReady {
// 		logrus.Info("Pods Unready: ", unready)
// 	}

// 	return clusterReady, nil
// }

// func IsPodReady(pod corev1.Pod) bool {
// 	succeeded := pod.Status.Phase == corev1.PodSucceeded

// 	for _, v := range pod.Status.Conditions {
// 		if v.Type == corev1.ContainersReady {
// 			running := v.Status == corev1.ConditionTrue
// 			return succeeded || running
// 		}
// 	}

// 	return false
// }
