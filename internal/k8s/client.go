package k8s

import (
	"context"
	"strings"

	helmv2beta1 "github.com/fluxcd/helm-controller/api/v2beta1"
	helmv2beta2 "github.com/fluxcd/helm-controller/api/v2beta2"
	kustomizev1 "github.com/fluxcd/kustomize-controller/api/v1"
	kustomizev1beta1 "github.com/fluxcd/kustomize-controller/api/v1beta1"
	kustomizev1beta2 "github.com/fluxcd/kustomize-controller/api/v1beta2"
	sourcev1 "github.com/fluxcd/source-controller/api/v1"
	sourcev1beta1 "github.com/fluxcd/source-controller/api/v1beta1"
	sourcev1beta2 "github.com/fluxcd/source-controller/api/v1beta2"

	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	ctlrconfig "sigs.k8s.io/controller-runtime/pkg/client/config"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"

	"github.com/pthomison/k3auto/internal/yaml"
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

	err = kustomizev1.AddToScheme(scheme)
	if err != nil {
		return nil, err
	}
	err = kustomizev1beta1.AddToScheme(scheme)
	if err != nil {
		return nil, err
	}
	err = kustomizev1beta2.AddToScheme(scheme)
	if err != nil {
		return nil, err
	}

	err = sourcev1.AddToScheme(scheme)
	if err != nil {
		return nil, err
	}
	err = sourcev1beta2.AddToScheme(scheme)
	if err != nil {
		return nil, err
	}

	err = sourcev1beta1.AddToScheme(scheme)
	if err != nil {
		return nil, err
	}

	err = helmv2beta1.AddToScheme(scheme)
	if err != nil {
		return nil, err
	}
	err = helmv2beta2.AddToScheme(scheme)
	if err != nil {
		return nil, err
	}

	err = apiextensionsv1.AddToScheme(scheme)
	if err != nil {
		return nil, err
	}
	err = apiextensionsv1beta1.AddToScheme(scheme)
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

func CreateManifests(ctx context.Context, k8sC client.Client, manifests string) error {

	rawObjs, err := yaml.YamlReadAndSplit(strings.NewReader(manifests))
	if err != nil {
		return err
	}

	decoder := NewDecoder()

	for i := range rawObjs {

		obj, _, err := decoder.Decode(rawObjs[i], nil, nil)
		if err != nil {
			return err
		}

		err = k8sC.Create(ctx, obj.(client.Object))
		if err != nil {
			return err
		}

	}

	return nil
}

func CreateObjects(ctx context.Context, k8sC client.Client, objs []runtime.Object) error {

	for _, obj := range objs {

		err := k8sC.Create(ctx, obj.(client.Object))
		if err != nil {
			return err
		}

	}

	return nil
}
