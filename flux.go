package main

import (
	"time"

	kustomizev1 "github.com/fluxcd/kustomize-controller/api/v1"
	"github.com/fluxcd/pkg/apis/meta"
	sourcev1 "github.com/fluxcd/source-controller/api/v1"
	"github.com/pthomison/utilkit"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	secret = corev1.Secret{
		ObjectMeta: v1.ObjectMeta{
			Name:      "flux-system",
			Namespace: "flux-system",
		},
		StringData: map[string]string{
			"identity": utilkit.MustRequestParameter(utilkit.RequestParameterInput{
				Name:   "/k3d-playground/deploy-key/identity",
				Region: "us-east-2",
			}),
			"identity.pub": utilkit.MustRequestParameter(utilkit.RequestParameterInput{
				Name:   "/k3d-playground/deploy-key/identity.pub",
				Region: "us-east-2",
			}),
			"known_hosts": utilkit.MustRequestParameter(utilkit.RequestParameterInput{
				Name:   "/k3d-playground/deploy-key/known_hosts",
				Region: "us-east-2",
			}),
		},
	}

	gitrepo = sourcev1.GitRepository{
		ObjectMeta: v1.ObjectMeta{
			Name:      "flux-system",
			Namespace: "flux-system",
		},
		Spec: sourcev1.GitRepositorySpec{
			Interval: v1.Duration{
				Duration: time.Minute * 1,
			},
			Reference: &sourcev1.GitRepositoryRef{
				Branch: "main",
			},
			SecretRef: &meta.LocalObjectReference{
				Name: "flux-system",
			},
			URL: "ssh://git@github.com/pthomison/flux-environments",
		},
	}

	kustomization = kustomizev1.Kustomization{
		ObjectMeta: v1.ObjectMeta{
			Name:      "flux-system",
			Namespace: "flux-system",
		},
		Spec: kustomizev1.KustomizationSpec{
			Interval: v1.Duration{
				Duration: time.Minute * 10,
			},
			Path:  "./clusters/k3auto",
			Prune: true,
			SourceRef: kustomizev1.CrossNamespaceSourceReference{
				Kind: "GitRepository",
				Name: "flux-system",
			},
		},
	}
)
