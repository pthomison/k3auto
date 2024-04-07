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
	Short: "Reinject deployments",
	Run:   k3AutoUpdate,
}

func k3AutoUpdate(cmd *cobra.Command, args []string) {
	ctx := cmd.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var err error

	if !MinimalFlag {

		logrus.Info("Injecting Default Deployments")
		err = Deploy(ctx, "default", defaults.DefaultDeploymentsFolder, "/", afero.FromIOFS{FS: defaults.DefaultDeployments})
		checkError(err)
		logrus.Info("Default Deployments Injected")

	}

	if DeploymentDirectoryFlag != "" {

		logrus.Info("Injecting Directory Deployments")
		err = Deploy(ctx, "deployments", DeploymentDirectoryFlag, BootstrapDirectoryFlag, afero.NewOsFs())
		checkError(err)

		logrus.Info("Directory Deployments Injected")

	}
}
