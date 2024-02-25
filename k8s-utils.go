package main

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func WaitForPodsReadInCluster(ctx context.Context, k8s client.Client) error {
	ready := false
	var err error
	for !ready {
		ready, err = ArePodsReadyInCluster(ctx, k8s)
		if err != nil {
			return err
		}

		logrus.Info("Waiting on cluster")
		time.Sleep(1 * time.Second)
	}
	return nil
}

func ArePodsReadyInCluster(ctx context.Context, k8s client.Client) (bool, error) {
	podList := corev1.PodList{}
	err := k8s.List(ctx, &podList)
	if err != nil {
		return false, err
	}

	clusterReady := true
	unready := []string{}

	for _, pod := range podList.Items {
		podReady := IsPodReady(pod)

		if !podReady {
			clusterReady = false
			unready = append(unready, pod.ObjectMeta.Name)
		}
	}

	logrus.Info("Pods Unready: ", unready)

	return clusterReady, nil
}

func IsPodReady(pod corev1.Pod) bool {
	succeeded := pod.Status.Phase == corev1.PodSucceeded

	for _, v := range pod.Status.Conditions {
		if v.Type == corev1.ContainersReady {
			running := v.Status == corev1.ConditionTrue
			return succeeded || running
		}
	}

	return false
}
