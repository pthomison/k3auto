package k3auto

import (
	"context"

	defaults "github.com/pthomison/k3auto/default"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

func Update(ctx context.Context, conf Config) error {
	if !conf.Minimal {

		logrus.Info("Injecting Default Deployments")
		err := Deploy(ctx, "default", defaults.DefaultDeploymentsFolder, "/", afero.FromIOFS{FS: defaults.DefaultDeployments})
		if err != nil {
			return err
		}
		logrus.Info("Default Deployments Injected")

	}

	if conf.DeploymentDirectory != "" {

		logrus.Info("Injecting Directory Deployments")
		err := Deploy(ctx, "deployments", conf.DeploymentDirectory, conf.BootstrapDirectory, afero.NewOsFs())
		if err != nil {
			return err
		}

		logrus.Info("Directory Deployments Injected")

	}
	return nil
}
