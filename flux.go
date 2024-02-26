package main

// import (
// 	"time"

// 	"github.com/pthomison/utilkit"

// 	"github.com/fluxcd/pkg/apis/meta"

// 	kustomizev1 "github.com/fluxcd/kustomize-controller/api/v1"
// 	sourcev1 "github.com/fluxcd/source-controller/api/v1"
// 	sourcev1beta2 "github.com/fluxcd/source-controller/api/v1beta2"

// 	corev1 "k8s.io/api/core/v1"
// 	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
// )

// var (
// 	secret = corev1.Secret{
// 		ObjectMeta: v1.ObjectMeta{
// 			Name:      "flux-system",
// 			Namespace: "flux-system",
// 		},
// 		StringData: map[string]string{
// 			"identity": utilkit.MustRequestParameter(utilkit.RequestParameterInput{
// 				Name:   "/k3d-playground/deploy-key/identity",
// 				Region: "us-east-2",
// 			}),
// 			"identity.pub": utilkit.MustRequestParameter(utilkit.RequestParameterInput{
// 				Name:   "/k3d-playground/deploy-key/identity.pub",
// 				Region: "us-east-2",
// 			}),
// 			"known_hosts": utilkit.MustRequestParameter(utilkit.RequestParameterInput{
// 				Name:   "/k3d-playground/deploy-key/known_hosts",
// 				Region: "us-east-2",
// 			}),
// 		},
// 	}

// 	gitrepo = sourcev1.GitRepository{
// 		ObjectMeta: v1.ObjectMeta{
// 			Name:      "flux-system",
// 			Namespace: "flux-system",
// 		},
// 		Spec: sourcev1.GitRepositorySpec{
// 			Interval: v1.Duration{
// 				Duration: time.Minute * 1,
// 			},
// 			Reference: &sourcev1.GitRepositoryRef{
// 				Branch: "main",
// 			},
// 			SecretRef: &meta.LocalObjectReference{
// 				Name: "flux-system",
// 			},
// 			URL: "ssh://git@github.com/pthomison/flux-environments",
// 		},
// 	}

// 	ocirepo = sourcev1beta2.OCIRepository{
// 		ObjectMeta: v1.ObjectMeta{
// 			Name:      "flux-system-oci",
// 			Namespace: "flux-system",
// 		},
// 		Spec: sourcev1beta2.OCIRepositorySpec{
// 			Interval: v1.Duration{
// 				Duration: time.Minute * 5,
// 			},
// 			URL: "oci://172.17.0.4:5000/hackstash",
// 			Reference: &sourcev1beta2.OCIRepositoryRef{
// 				Tag: "latest",
// 			},
// 			Insecure: true,
// 		},
// 	}

// 	kustomizationOCI = kustomizev1.Kustomization{
// 		ObjectMeta: v1.ObjectMeta{
// 			Name:      "flux-system",
// 			Namespace: "flux-system",
// 		},
// 		Spec: kustomizev1.KustomizationSpec{
// 			Interval: v1.Duration{
// 				Duration: time.Minute * 10,
// 			},
// 			Path:  "/",
// 			Prune: true,
// 			SourceRef: kustomizev1.CrossNamespaceSourceReference{
// 				Kind:      "OCIRepository",
// 				Name:      "flux-system-oci",
// 				Namespace: "flux-system",
// 			},
// 		},
// 	}

// 	kustomization = kustomizev1.Kustomization{
// 		ObjectMeta: v1.ObjectMeta{
// 			Name:      "flux-system",
// 			Namespace: "flux-system",
// 		},
// 		Spec: kustomizev1.KustomizationSpec{
// 			Interval: v1.Duration{
// 				Duration: time.Minute * 10,
// 			},
// 			Path:  "./clusters/k3d-playground",
// 			Prune: true,
// 			SourceRef: kustomizev1.CrossNamespaceSourceReference{
// 				Kind: "GitRepository",
// 				Name: "flux-system",
// 			},
// 		},
// 	}
// )
