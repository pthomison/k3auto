package k3d

import (
	"context"
	"embed"
	"errors"

	"github.com/spf13/afero"
	"github.com/spf13/viper"

	k3dconfig "github.com/k3d-io/k3d/v5/pkg/config"
	k3dv1alpha5 "github.com/k3d-io/k3d/v5/pkg/config/v1alpha5"
	k3druntimes "github.com/k3d-io/k3d/v5/pkg/runtimes"
)

func ParseConfigFile(confLocation string, embedFs *embed.FS) (*k3dv1alpha5.SimpleConfig, error) {
	config := viper.New()

	if embedFs != nil {
		config.SetFs(afero.FromIOFS{FS: *embedFs})
	}

	config.SetConfigFile(confLocation)

	if err := config.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, err
		}
		return nil, err
	}

	cfg, err := k3dconfig.FromViper(config)
	if err != nil {
		return nil, err
	}

	versionConf, ok := cfg.(k3dv1alpha5.SimpleConfig)
	if !ok {
		return nil, errors.New("failed to cast config file to v1alpha5.SimpleConfig")
	}

	return &versionConf, nil
}

func LoadClusterConfig(ctx context.Context, rt k3druntimes.Runtime, cfg *k3dv1alpha5.SimpleConfig) (*k3dv1alpha5.ClusterConfig, error) {
	err := k3dconfig.ProcessSimpleConfig(cfg)
	if err != nil {
		return nil, err
	}

	clusterConfig, err := k3dconfig.TransformSimpleToClusterConfig(ctx, rt, *cfg, "")
	if err != nil {
		return nil, err
	}

	err = k3dconfig.ValidateClusterConfig(ctx, rt, *clusterConfig)
	if err != nil {
		return nil, err
	}

	return clusterConfig, nil
}
