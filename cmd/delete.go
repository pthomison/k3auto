package cmd

import (
	"context"

	k3druntimes "github.com/k3d-io/k3d/v5/pkg/runtimes"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/pthomison/k3auto/internal/k3d"
)

var DeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete an existing cluster",
	Run:   k3AutoDelete,
}

func k3AutoDelete(cmd *cobra.Command, args []string) {
	ctx := cmd.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	clusterConfig, err := parseConfigFile(ClusterConfigFileFlag)
	checkError(err)
	logrus.Info("K3D Config File Loaded: ", ClusterConfigFileFlag)

	logrus.Info("Deleting Cluster")
	err = k3d.DeleteCluster(ctx, clusterConfig, k3druntimes.Docker)
	checkError(err)
	logrus.Info("Cluster Deleted")
}
