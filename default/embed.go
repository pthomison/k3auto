package defaults

import (
	"embed"
)

const (
	K3dConfigLocation = "k3d-config.yaml"
)

var (
	//go:embed k3d-config.yaml
	K3dConfig embed.FS

	//go:embed deployments
	DefaultDeployments embed.FS
)
