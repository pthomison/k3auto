package main

import (
	"fmt"
	"time"

	configtypes "github.com/k3d-io/k3d/v5/pkg/config/types"
	v1alpha5 "github.com/k3d-io/k3d/v5/pkg/config/v1alpha5"
	"github.com/k3d-io/k3d/v5/pkg/types"
)

var (
	clusterSimpleCfg = &v1alpha5.SimpleConfig{
		ObjectMeta: configtypes.ObjectMeta{
			Name: "foobar",
		},
		Servers: 1,
		Image:   fmt.Sprintf("%s:%s", types.DefaultK3sImageRepo, "v1.29.1-k3s1"),
		ExposeAPI: v1alpha5.SimpleExposureOpts{
			HostPort: "6443",
		},
		Network: "bridge",
		Options: v1alpha5.SimpleConfigOptions{
			K3dOptions: v1alpha5.SimpleConfigOptionsK3d{
				Wait:                true,
				Timeout:             3 * time.Minute,
				DisableLoadbalancer: true,
			},
			K3sOptions: v1alpha5.SimpleConfigOptionsK3s{
				ExtraArgs: []v1alpha5.K3sArgWithNodeFilters{
					{
						Arg:         "--disable=servicelb,traefik,metrics-server",
						NodeFilters: []string{"server:*"},
					},
					{
						Arg:         "--disable-network-policy",
						NodeFilters: []string{"server:*"},
					},
				},
			},
			KubeconfigOptions: v1alpha5.SimpleConfigOptionsKubeconfig{
				UpdateDefaultKubeconfig: true,
				SwitchCurrentContext:    true,
			},
			Runtime: v1alpha5.SimpleConfigOptionsRuntime{},
		},
	}
)
