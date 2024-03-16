package k3d

import (
	"context"

	k3dcluster "github.com/k3d-io/k3d/v5/pkg/client"
	k3dconfig "github.com/k3d-io/k3d/v5/pkg/config"
	v1alpha5 "github.com/k3d-io/k3d/v5/pkg/config/v1alpha5"
	"github.com/k3d-io/k3d/v5/pkg/runtimes"
)

func DeployCluster(ctx context.Context, cfg *v1alpha5.SimpleConfig, rt runtimes.Runtime) error {
	err := k3dconfig.ProcessSimpleConfig(cfg)
	if err != nil {
		return err
	}

	clusterConfig, err := k3dconfig.TransformSimpleToClusterConfig(ctx, rt, *cfg)
	if err != nil {
		return err
	}

	err = k3dconfig.ValidateClusterConfig(ctx, runtimes.SelectedRuntime, *clusterConfig)
	if err != nil {
		return err
	}

	err = k3dcluster.ClusterRun(ctx, rt, clusterConfig)
	if err != nil {
		return err
	}

	_, err = k3dcluster.KubeconfigGetWrite(ctx, runtimes.SelectedRuntime,
		&clusterConfig.Cluster,
		"",
		&k3dcluster.WriteKubeConfigOptions{
			UpdateExisting:       true,
			OverwriteExisting:    true,
			UpdateCurrentContext: cfg.Options.KubeconfigOptions.SwitchCurrentContext,
		},
	)
	if err != nil {
		return err
	}

	return nil
}
