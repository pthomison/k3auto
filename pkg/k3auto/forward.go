package k3auto

import (
	"context"

	"github.com/pthomison/k3auto/internal/k8s"
)

func PortForward(ctx context.Context, name string, namespace string, port int) (chan struct{}, error) {
	return k8s.PortForward(ctx, name, namespace, port)
}
