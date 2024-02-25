package main

import (
	"context"

	k3cluster "github.com/k3d-io/k3d/v5/pkg/client"
	"github.com/k3d-io/k3d/v5/pkg/config"
	v1alpha5 "github.com/k3d-io/k3d/v5/pkg/config/v1alpha5"
	"github.com/k3d-io/k3d/v5/pkg/runtimes"
)

func DeployCluster(ctx context.Context, cfg *v1alpha5.SimpleConfig) error {
	err := config.ProcessSimpleConfig(clusterSimpleCfg)
	if err != nil {
		return err
	}

	clusterConfig, err := config.TransformSimpleToClusterConfig(ctx, rt, *clusterSimpleCfg)
	if err != nil {
		return err
	}

	err = config.ValidateClusterConfig(ctx, runtimes.SelectedRuntime, *clusterConfig)
	if err != nil {
		return err
	}

	err = k3cluster.ClusterRun(ctx, rt, clusterConfig)
	if err != nil {
		return err
	}

	if _, err := k3cluster.KubeconfigGetWrite(ctx, runtimes.SelectedRuntime,
		&clusterConfig.Cluster,
		"",
		&k3cluster.WriteKubeConfigOptions{
			UpdateExisting:       true,
			OverwriteExisting:    true,
			UpdateCurrentContext: clusterSimpleCfg.Options.KubeconfigOptions.SwitchCurrentContext,
		},
	); err != nil {
		return err
	}

	return nil
}
