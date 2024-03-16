package flux

import (
	"context"
	"time"

	kustomizev1 "github.com/fluxcd/kustomize-controller/api/v1"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apimachinerytypes "k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func WaitForKustomization(ctx context.Context, k8s client.Client, desiredDep v1.ObjectMeta) error {
	for {
		k := kustomizev1.Kustomization{}
		err := k8s.Get(ctx, apimachinerytypes.NamespacedName{
			Name:      desiredDep.Name,
			Namespace: desiredDep.Namespace,
		}, &k)
		if err != nil {
			return err
		}

		for _, cond := range k.Status.Conditions {
			if cond.Type == "Ready" {
				if cond.Status == "True" {
					return nil
				}
			}
		}
		time.Sleep(time.Second * 1)
	}
}
