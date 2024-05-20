package flux

import (
	"github.com/fluxcd/flux2/v2/pkg/manifestgen"
	"github.com/fluxcd/flux2/v2/pkg/manifestgen/install"
)

func GenerateManifests(version string) (*manifestgen.Manifest, error) {
	// Generate Flux Controller Manifests
	genOps := install.MakeDefaultOptions()
	genOps.NetworkPolicy = false
	genOps.Version = version
	manifests, err := install.Generate(genOps, "")
	if err != nil {
		return nil, err
	}

	return manifests, nil
}
