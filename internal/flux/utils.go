package flux

import (
	"context"
	"fmt"
	"time"

	kustomizev1 "github.com/fluxcd/kustomize-controller/api/v1"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apimachinerytypes "k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	sourcev1beta2 "github.com/fluxcd/source-controller/api/v1beta2"
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

func NewOCIKustomization(name string, image string, tag string) (sourcev1beta2.OCIRepository, kustomizev1.Kustomization) {

	ocirepo := sourcev1beta2.OCIRepository{
		ObjectMeta: v1.ObjectMeta{
			Name:      name,
			Namespace: "flux-system",
		},
		Spec: sourcev1beta2.OCIRepositorySpec{
			Interval: v1.Duration{
				Duration: time.Minute * 5,
			},
			URL: fmt.Sprintf("oci://k3auto-registry.localhost:5000/%v", image),
			Reference: &sourcev1beta2.OCIRepositoryRef{
				Tag: tag,
			},
			Insecure: true,
		},
	}

	kustomizationOCI := kustomizev1.Kustomization{
		ObjectMeta: v1.ObjectMeta{
			Name:      name,
			Namespace: "flux-system",
		},
		Spec: kustomizev1.KustomizationSpec{
			Interval: v1.Duration{
				Duration: time.Minute * 10,
			},
			Path:  "/",
			Prune: true,
			SourceRef: kustomizev1.CrossNamespaceSourceReference{
				Kind:      "OCIRepository",
				Name:      name,
				Namespace: "flux-system",
			},
		},
	}

	return ocirepo, kustomizationOCI
}
