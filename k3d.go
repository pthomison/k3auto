package main

import (
	"context"

	k3dcluster "github.com/k3d-io/k3d/v5/pkg/client"
	k3dconfig "github.com/k3d-io/k3d/v5/pkg/config"
	v1alpha5 "github.com/k3d-io/k3d/v5/pkg/config/v1alpha5"
	"github.com/k3d-io/k3d/v5/pkg/runtimes"
)

func DeployCluster(ctx context.Context, cfg *v1alpha5.SimpleConfig) error {
	err := k3dconfig.ProcessSimpleConfig(clusterSimpleCfg)
	if err != nil {
		return err
	}

	clusterConfig, err := k3dconfig.TransformSimpleToClusterConfig(ctx, rt, *clusterSimpleCfg)
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
			UpdateCurrentContext: clusterSimpleCfg.Options.KubeconfigOptions.SwitchCurrentContext,
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// func GetCluster(ctx context.Context, cfg *v1alpha5.SimpleConfig) error {
// 	clusters, err := k3dcluster.ClusterList(ctx, rt)
// 	if err != nil {
// 		return err
// 	}

// 	registry, err := k3dcluster.RegistryGet(ctx, rt, "k3auto-registry")
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
