package k3auto

import (
	"context"

	k3druntimes "github.com/k3d-io/k3d/v5/pkg/runtimes"
	"github.com/pthomison/k3auto/internal/k3d"
	"github.com/sirupsen/logrus"
)

func Delete(ctx context.Context, conf Config) error {

	clusterConfig, err := ParseK3dConfigFile(conf.ClusterConfigFile)
	if err != nil {
		return err
	}
	logrus.Info("K3D Config File Loaded: ", conf.ClusterConfigFile)

	logrus.Info("Deleting Cluster")
	err = k3d.DeleteCluster(ctx, clusterConfig, k3druntimes.Docker)
	if err != nil {
		return err
	}
	logrus.Info("Cluster Deleted")

	return nil
}
