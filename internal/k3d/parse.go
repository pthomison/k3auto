package k3d

import (
	"errors"

	"github.com/spf13/viper"

	k3dconfig "github.com/k3d-io/k3d/v5/pkg/config"
	v1alpha5 "github.com/k3d-io/k3d/v5/pkg/config/v1alpha5"
)

func ParseConfigFile(confLocation string) (*v1alpha5.SimpleConfig, error) {
	config := viper.New()
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

	versionConf, ok := cfg.(v1alpha5.SimpleConfig)
	if !ok {
		return nil, errors.New("failed to cast config file to v1alpha5.SimpleConfig")
	}

	return &versionConf, nil
}
