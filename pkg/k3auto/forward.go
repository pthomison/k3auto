package k3auto

import (
	"context"

	"github.com/pthomison/k3auto/internal/k8s"
	apimachinerytypes "k8s.io/apimachinery/pkg/types"
)

func PortForward(ctx context.Context, name apimachinerytypes.NamespacedName, port int) (chan struct{}, error) {
	kcfg, _, err := k8s.NewClient()
	if err != nil {
		return nil, err
	}

	return k8s.PortForward(ctx, kcfg, name, port)
}

func PortForwardService(ctx context.Context, name apimachinerytypes.NamespacedName, port int) (chan struct{}, error) {
	kcfg, k8sC, err := k8s.NewClient()
	if err != nil {
		return nil, err
	}

	return k8s.PortForwardService(ctx, k8sC, kcfg, name, port)
}

func PortForwardDeployment(ctx context.Context, name apimachinerytypes.NamespacedName, port int) (chan struct{}, error) {
	kcfg, k8sC, err := k8s.NewClient()
	if err != nil {
		return nil, err
	}

	return k8s.PortForwardDeployment(ctx, k8sC, kcfg, name, port)
}
