package k8s

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"

	kustomizev1 "github.com/fluxcd/kustomize-controller/api/v1"
	kustomizev1beta1 "github.com/fluxcd/kustomize-controller/api/v1beta1"
	kustomizev1beta2 "github.com/fluxcd/kustomize-controller/api/v1beta2"

	sourcev1 "github.com/fluxcd/source-controller/api/v1"
	sourcev1beta1 "github.com/fluxcd/source-controller/api/v1beta1"
	sourcev1beta2 "github.com/fluxcd/source-controller/api/v1beta2"

	helmv2beta1 "github.com/fluxcd/helm-controller/api/v2beta1"
	helmv2beta2 "github.com/fluxcd/helm-controller/api/v2beta2"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
)

func NewDecoder() runtime.Decoder {
	apiextensionsv1.AddToScheme(scheme.Scheme)
	apiextensionsv1beta1.AddToScheme(scheme.Scheme)

	kustomizev1.AddToScheme(scheme.Scheme)
	kustomizev1beta1.AddToScheme(scheme.Scheme)
	kustomizev1beta2.AddToScheme(scheme.Scheme)

	sourcev1.AddToScheme(scheme.Scheme)
	sourcev1beta2.AddToScheme(scheme.Scheme)
	sourcev1beta1.AddToScheme(scheme.Scheme)

	helmv2beta1.AddToScheme(scheme.Scheme)
	helmv2beta2.AddToScheme(scheme.Scheme)

	decoder := scheme.Codecs.UniversalDeserializer()

	return decoder
}
