package k3auto

import (
	"fmt"

	"github.com/davecgh/go-spew/spew"
	k3dv1alpha5 "github.com/k3d-io/k3d/v5/pkg/config/v1alpha5"
	defaults "github.com/pthomison/k3auto/default"
	"github.com/pthomison/k3auto/internal/k3d"
	"github.com/spf13/afero"
)

type Config struct {
	Minimal bool

	DeploymentDirectory string

	DeploymentFilesystem afero.Fs

	BootstrapDirectory string
	SecretFile         string
	ClusterConfigFile  string
}

func ParseK3dConfigFile(configPath string) (*k3dv1alpha5.SimpleConfig, error) {
	var clusterConfig *k3dv1alpha5.SimpleConfig
	var err error

	if configPath != "" {
		clusterConfig, err = k3d.ParseConfigFile(configPath, nil)
		if err != nil {
			return nil, err
		}

	} else {
		clusterConfig, err = k3d.ParseConfigFile(defaults.K3dConfigLocation, &defaults.K3dConfig)
		if err != nil {
			return nil, err
		}
		clusterConfig.Image = fmt.Sprintf("docker.io/rancher/k3s:v1.29.4-k3s1")
		spew.Dump(clusterConfig.Image)
	}

	return clusterConfig, nil
}
