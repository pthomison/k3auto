package defaults

import (
	"embed"
)

const (
	K3dConfigLocation = "k3d-config.yaml"

	DefaultDeploymentsFolder = "deployments"

	// RegistryDeploymentFolder = "registry"
)

var (
	//go:embed k3d-config.yaml
	K3dConfig embed.FS

	//go:embed deployments/*
	DefaultDeployments embed.FS

	//go:embed registry/registry.yaml
	RegistryDeployment string
)
