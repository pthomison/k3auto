package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Implementation Pending; Update the deployments in a cluster",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Version: %v\nCommit: %v\nDate: %v\n", version, commit, date)
	},
}

// func k3AutoUpdate(cmd *cobra.Command, args []string) {
// 	ctx := cmd.Context()
// 	if ctx == nil {
// 		ctx = context.Background()
// 	}

// 	var err error

// 	// clusterConfig, err := parseConfigFile(ClusterConfigFileFlag)
// 	// checkError(err)
// 	// logrus.Info("K3D Config File Loaded: ", ClusterConfigFileFlag)

// 	if !MinimalFlag {

// 		logrus.Info("Injecting Default Deployments")
// 		err = Deploy(ctx, "default", defaults.DefaultDeploymentsFolder, afero.FromIOFS{FS: defaults.DefaultDeployments})
// 		checkError(err)
// 		logrus.Info("Default Deployments Injected")

// 	}

// 	if DeploymentDirectoryFlag != "" {

// 		logrus.Info("Injecting Directory Deployments")
// 		err = Deploy(ctx, "deployments", DeploymentDirectoryFlag, afero.NewOsFs())
// 		checkError(err)

// 		logrus.Info("Directory Deployments Injected")

// 	}
// }
