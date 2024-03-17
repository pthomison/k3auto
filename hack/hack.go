package hack

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/scheme"

	kustomizev1 "github.com/fluxcd/kustomize-controller/api/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
)

func deserialize(data []byte) (runtime.Object, *schema.GroupVersionKind, error) {
	apiextensionsv1.AddToScheme(scheme.Scheme)
	apiextensionsv1beta1.AddToScheme(scheme.Scheme)
	kustomizev1.AddToScheme(scheme.Scheme)

	decoder := scheme.Codecs.UniversalDeserializer()

	runtimeObject, groupVersionKind, err := decoder.Decode(data, nil, nil)
	if err != nil {
		return nil, nil, err
	}

	return runtimeObject, groupVersionKind, nil
}
