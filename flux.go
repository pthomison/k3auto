package main

import (
	"time"

	kustomizev1 "github.com/fluxcd/kustomize-controller/api/v1"
	sourcev1beta2 "github.com/fluxcd/source-controller/api/v1beta2"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	ocirepo = sourcev1beta2.OCIRepository{
		ObjectMeta: v1.ObjectMeta{
			Name:      "flux-system-oci",
			Namespace: "flux-system",
		},
		Spec: sourcev1beta2.OCIRepositorySpec{
			Interval: v1.Duration{
				Duration: time.Minute * 5,
			},
			URL: "oci://172.17.0.4:5000/hackstash",
			Reference: &sourcev1beta2.OCIRepositoryRef{
				Tag: "latest",
			},
			Insecure: true,
		},
	}

	kustomizationOCI = kustomizev1.Kustomization{
		ObjectMeta: v1.ObjectMeta{
			Name:      "flux-system",
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
				Name:      "flux-system-oci",
				Namespace: "flux-system",
			},
		},
	}
)
