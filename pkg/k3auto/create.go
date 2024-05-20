package k3auto

import (
	"context"
	"os"

	k3dv1alpha5 "github.com/k3d-io/k3d/v5/pkg/config/v1alpha5"
	k3druntimes "github.com/k3d-io/k3d/v5/pkg/runtimes"
	defaults "github.com/pthomison/k3auto/default"
	"github.com/pthomison/k3auto/internal/flux"
	"github.com/pthomison/k3auto/internal/k3d"
	"github.com/pthomison/k3auto/internal/k8s"
	"github.com/pthomison/k3auto/internal/secrets"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func Create(ctx context.Context, conf Config) (ctrlclient.Client, error) {
	clusterConfig, err := ParseK3dConfigFile(conf.ClusterConfigFile)
	if err != nil {
		return nil, err
	}
	logrus.Info("K3D Config File Loaded: ", conf.ClusterConfigFile)

	logrus.Info("Initializing Cluster")
	k8sC, err := InitializeCluster(ctx, clusterConfig, k3druntimes.Docker)
	if err != nil {
		return nil, err
	}
	logrus.Info("Cluster Initialized")

	logrus.Info("Injecting Flux Controllers")
	err = InjectFluxControllers(ctx, k8sC, conf.FluxVersion)
	if err != nil {
		return nil, err
	}
	logrus.Info("Flux Controllers Injected")

	logrus.Info("Injecting Registry")
	err = InjectRegistry(ctx, k8sC)
	if err != nil {
		return nil, err
	}
	logrus.Info("Registry Injected")

	err = k8s.WaitForDeployment(ctx, k8sC, v1.ObjectMeta{
		Name:      "docker-registry",
		Namespace: "docker-registry",
	})
	if err != nil {
		return nil, err
	}

	if conf.SecretFile != "" {
		err = InjectSecrets(ctx, k8sC, conf.SecretFile)
		if err != nil {
			return nil, err
		}
	}

	if !conf.Minimal {
		logrus.Info("Injecting Default Deployments")
		err = Deploy(ctx, "default", defaults.DefaultDeploymentsFolder, "/", afero.FromIOFS{FS: defaults.DefaultDeployments})
		if err != nil {
			return nil, err
		}
		logrus.Info("Default Deployments Injected")
	}

	if conf.DeploymentDirectory != "" {
		logrus.Info("Injecting Directory Deployments")
		err = Deploy(ctx, "deployments", conf.DeploymentDirectory, conf.BootstrapDirectory, conf.DeploymentFilesystem)
		if err != nil {
			return nil, err
		}

		logrus.Info("Directory Deployments Injected")
	}

	return k8sC, nil
}

func InitializeCluster(ctx context.Context, config *k3dv1alpha5.SimpleConfig, runtime k3druntimes.Runtime) (ctrlclient.Client, error) {
	// Deploy the cluster defined in cluster.go
	err := k3d.DeployCluster(ctx, config, runtime)
	if err != nil {
		return nil, err
	}

	// Generate a k8s client from standard kubeconfig
	_, k8sC, err := k8s.NewClient()
	if err != nil {
		return nil, err
	}

	// Wait for the base cluster deployments to be ready
	k8s.WaitForDeployment(ctx, k8sC, v1.ObjectMeta{
		Name:      "coredns",
		Namespace: "kube-system",
	})

	return k8sC, err
}

func InjectFluxControllers(ctx context.Context, k8sC ctrlclient.Client, version string) error {
	fluxManifests, err := flux.GenerateManifests(version)
	if err != nil {
		return err
	}

	err = k8s.CreateManifests(ctx, k8sC, fluxManifests.Content)
	if err != nil {
		return err
	}

	return nil
}

func InjectRegistry(ctx context.Context, k8sC ctrlclient.Client) error {
	err := k8s.CreateManifests(ctx, k8sC, defaults.RegistryDeployment)
	if err != nil {
		return err
	}

	return nil
}

func InjectSecrets(ctx context.Context, k8sC ctrlclient.Client, secretFileLocation string) error {
	f, err := os.Open(secretFileLocation)
	if err != nil {
		return err
	}
	defer f.Close()

	conf, err := secrets.LoadConfigFile(f)
	if err != nil {
		return err
	}

	err = secrets.InjectSecrets(ctx, k8sC, conf)
	if err != nil {
		return err
	}

	return nil
}
