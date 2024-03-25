package cmd

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var UpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Implementation Pending; Update the deployments in a cluster",
	Run:   k3AutoUpdate,
}

func k3AutoUpdate(cmd *cobra.Command, args []string) {
	ctx := cmd.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	clusterConfig, err := parseConfigFile(ClusterConfigFileFlag)
	checkError(err)
	logrus.Info("K3D Config File Loaded: ", ClusterConfigFileFlag)

	_ = clusterConfig
	logrus.Info("Implementation Pending")
}
