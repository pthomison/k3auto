package flux

import (
	"github.com/fluxcd/flux2/v2/pkg/manifestgen"
	"github.com/fluxcd/flux2/v2/pkg/manifestgen/install"
)

func GenerateManifests() (*manifestgen.Manifest, error) {
	// Generate Flux Controller Manifests
	genOps := install.MakeDefaultOptions()
	genOps.NetworkPolicy = false
	genOps.Version = "v2.3.0"
	manifests, err := install.Generate(genOps, "")
	if err != nil {
		return nil, err
	}

	// err = os.WriteFile("/tmp/test-manifests.yaml", []byte(manifests.Content), 0644)
	// if err != nil {
	// 	return nil, err
	// }

	return manifests, nil
}
