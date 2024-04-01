package cmd

import (
	"context"

	defaults "github.com/pthomison/k3auto/default"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
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

	var err error

	// clusterConfig, err := parseConfigFile(ClusterConfigFileFlag)
	// checkError(err)
	// logrus.Info("K3D Config File Loaded: ", ClusterConfigFileFlag)

	if !MinimalFlag {

		logrus.Info("Injecting Default Deployments")
		err = k3autoDeploy(ctx, "default", defaults.DefaultDeploymentsFolder, afero.FromIOFS{FS: defaults.DefaultDeployments})
		checkError(err)
		logrus.Info("Default Deployments Injected")

	}

	if DeploymentDirectoryFlag != "" {

		logrus.Info("Injecting Directory Deployments")
		err = k3autoDeploy(ctx, "deployments", DeploymentDirectoryFlag, afero.NewOsFs())
		checkError(err)

		logrus.Info("Directory Deployments Injected")

	}
}
