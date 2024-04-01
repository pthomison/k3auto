package cmd

import (
	"context"

	k3dv1alpha5 "github.com/k3d-io/k3d/v5/pkg/config/v1alpha5"
	k3druntimes "github.com/k3d-io/k3d/v5/pkg/runtimes"
	defaults "github.com/pthomison/k3auto/default"
	"github.com/pthomison/k3auto/internal/flux"
	"github.com/pthomison/k3auto/internal/k3d"
	"github.com/pthomison/k3auto/internal/k8s"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new K3D Cluster and inject flux controllers & deployments",
	Run:   k3AutoCreate,
}

func init() {
	opts := zap.Options{
		Development: true,
	}
	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))
}

func initializeCluster(ctx context.Context, config *k3dv1alpha5.SimpleConfig, runtime k3druntimes.Runtime) (ctrlclient.Client, error) {
	// Deploy the cluster defined in cluster.go
	err := k3d.DeployCluster(ctx, config, runtime)
	if err != nil {
		return nil, err
	}

	// Generate a k8s client from standard kubeconfig
	k8sC, err := k8s.NewClient()
	checkError(err)

	// Wait for the base cluster deployments to be ready
	k8s.WaitForDeployment(ctx, k8sC, v1.ObjectMeta{
		Name:      "coredns",
		Namespace: "kube-system",
	})

	return k8sC, err
}

func injectFluxControllers(ctx context.Context) error {
	fluxManifests, err := flux.GenerateManifests()
	if err != nil {
		return err
	}

	k8sC, err := k8s.NewClient()
	if err != nil {
		return err
	}

	err = k8s.CreateManifests(ctx, k8sC, fluxManifests.Content)
	if err != nil {
		return err
	}

	return nil
}

func k3AutoCreate(cmd *cobra.Command, args []string) {
	ctx := cmd.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	clusterConfig, err := parseConfigFile(ClusterConfigFileFlag)
	checkError(err)
	logrus.Info("K3D Config File Loaded: ", ClusterConfigFileFlag)

	logrus.Info("Initializing Cluster")
	_, err = initializeCluster(ctx, clusterConfig, k3druntimes.Docker)
	checkError(err)
	logrus.Info("Cluster Initialized")

	logrus.Info("Injecting Flux Controllers")
	err = injectFluxControllers(ctx)
	checkError(err)
	logrus.Info("Flux Controllers Injected")

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
