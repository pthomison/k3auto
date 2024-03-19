package k8s

import (
	kustomizev1 "github.com/fluxcd/kustomize-controller/api/v1"
	sourcev1 "github.com/fluxcd/source-controller/api/v1"
	sourcev1beta2 "github.com/fluxcd/source-controller/api/v1beta2"

	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	ctlrconfig "sigs.k8s.io/controller-runtime/pkg/client/config"
)

func NewClient() (client.Client, error) {
	kcfg, err := ctlrconfig.GetConfig()
	if err != nil {
		return nil, err
	}

	scheme := runtime.NewScheme()
	err = clientgoscheme.AddToScheme(scheme)
	if err != nil {
		return nil, err
	}
	err = sourcev1.AddToScheme(scheme)
	if err != nil {
		return nil, err
	}
	err = kustomizev1.AddToScheme(scheme)
	if err != nil {
		return nil, err
	}
	err = sourcev1beta2.AddToScheme(scheme)
	if err != nil {
		return nil, err
	}

	k8s, err := client.New(kcfg, client.Options{
		Scheme: scheme,
	})
	if err != nil {
		return nil, err
	}

	return k8s, err
}
