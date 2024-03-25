package k3d

import (
	"context"

	k3dcluster "github.com/k3d-io/k3d/v5/pkg/client"
	v1alpha5 "github.com/k3d-io/k3d/v5/pkg/config/v1alpha5"
	"github.com/k3d-io/k3d/v5/pkg/runtimes"
	k3dtypes "github.com/k3d-io/k3d/v5/pkg/types"
)

func DeployCluster(ctx context.Context, cfg *v1alpha5.SimpleConfig, rt runtimes.Runtime) error {
	clusterConfig, err := LoadClusterConfig(ctx, rt, cfg)
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

func DeleteCluster(ctx context.Context, cfg *v1alpha5.SimpleConfig, rt runtimes.Runtime) error {

	cluster, err := k3dcluster.ClusterGet(ctx, rt, &k3dtypes.Cluster{
		Name: cfg.Name,
	})
	if err != nil {
		return err
	}

	err = k3dcluster.ClusterDelete(ctx, rt, cluster, k3dtypes.ClusterDeleteOpts{})
	if err != nil {
		return err
	}

	err = k3dcluster.KubeconfigRemoveClusterFromDefaultConfig(ctx, cluster)
	if err != nil {
		return err
	}

	return nil
}
