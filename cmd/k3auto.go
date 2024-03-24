package cmd

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"io"
	"os"
	"path"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	ctrl "sigs.k8s.io/controller-runtime"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"

	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	k3dv1alpha5 "github.com/k3d-io/k3d/v5/pkg/config/v1alpha5"
	k3druntimes "github.com/k3d-io/k3d/v5/pkg/runtimes"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	defaults "github.com/pthomison/k3auto/default"
	"github.com/pthomison/k3auto/internal/flux"
	"github.com/pthomison/k3auto/internal/k3d"
	"github.com/pthomison/k3auto/internal/k8s"
)

var K3AutoCmd = &cobra.Command{
	Use:   "k3auto",
	Short: "k3auto is a local kubernetes cluster orchestrator powered by k3d and flux",
	Run:   k3AutoRun,
}

var (
	ClusterConfigFileFlag   string
	DeploymentDirectoryFlag string
	MinimalFlag             bool
)

func init() {
	K3AutoCmd.PersistentFlags().StringVarP(&ClusterConfigFileFlag, "cluster-config", "c", "", "Override Cluster Config File")
	K3AutoCmd.PersistentFlags().StringVarP(&DeploymentDirectoryFlag, "deployment-directory", "d", "", "Deployment Directory")
	K3AutoCmd.PersistentFlags().BoolVarP(&MinimalFlag, "minimal", "m", false, "Only deploy the k3d cluster & flux controllers")

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

func parseConfigFile(configPath string) (*k3dv1alpha5.SimpleConfig, error) {
	var clusterConfig *k3dv1alpha5.SimpleConfig
	var err error

	if configPath != "" {
		clusterConfig, err = k3d.ParseConfigFile(configPath, nil)
		if err != nil {
			return nil, err
		}
	} else {
		clusterConfig, err = k3d.ParseConfigFile(defaults.K3dConfigLocation, &defaults.K3dConfig)
		if err != nil {
			return nil, err
		}
	}

	return clusterConfig, nil
}

func yamlReadAndSplit(reader io.Reader) ([][]byte, error) {
	fb, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	// objs := bytes.Split(fb, []byte("---"))
	var objs [][]byte

	for _, obj := range bytes.Split(fb, []byte("---")) {
		if len(obj) != 0 {
			objs = append(objs, obj)
		}
	}

	return objs, err
}

func injectFluxControllers() error {
	tmpDirLoc, err := os.MkdirTemp("", "k3auto-flux-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDirLoc)

	// Generate Flux Controller Manifests
	fluxManifests, err := flux.GenerateManifests()
	if err != nil {
		return err
	}
	fluxManifestsPath := path.Join(tmpDirLoc, "flux-manifests.yaml")
	os.WriteFile(fluxManifestsPath, []byte(fluxManifests.Content), 0644)

	err = k8s.KubectlApply(fluxManifestsPath)
	if err != nil {
		return err
	}

	return nil
}

func k3AutoRun(cmd *cobra.Command, args []string) {
	ctx := cmd.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	clusterConfig, err := parseConfigFile(ClusterConfigFileFlag)
	checkError(err)
	logrus.Info("K3D Config File Loaded: ", ClusterConfigFileFlag)

	logrus.Info("Initializing Cluster")
	k8sC, err := initializeCluster(ctx, clusterConfig, k3druntimes.Docker)
	checkError(err)
	logrus.Info("Cluster Initialized")

	logrus.Info("Injecting Flux Controllers")
	err = injectFluxControllers()
	checkError(err)
	logrus.Info("Flux Controllers Injected")

	if !MinimalFlag {

		logrus.Info("Injecting Default Deployments")
		deploymentFiles, err := defaults.DefaultDeployments.ReadDir(defaults.DefaultDeploymentsFolder)
		checkError(err)
		for _, v := range deploymentFiles {
			f, err := defaults.DefaultDeployments.Open(fmt.Sprintf("%v/%v", defaults.DefaultDeploymentsFolder, v.Name()))
			checkError(err)
			defer f.Close()

			fileObjects, err := yamlReadAndSplit(f)
			checkError(err)

			for _, obj := range fileObjects {
				obj, objType, err := k8s.ParseManifest(obj)
				checkError(err)

				logrus.Info("Deploying: ", objType)

				err = k8sC.Create(ctx, obj.(ctrlclient.Object))
				checkError(err)
			}
		}
		logrus.Info("Default Deployments Injected")

	}

	if DeploymentDirectoryFlag != "" {

		logrus.Info("Injecting Directory Deployments")
		deploymentFiles, err := os.ReadDir(DeploymentDirectoryFlag)
		checkError(err)
		for _, v := range deploymentFiles {
			f, err := os.Open(fmt.Sprintf("%v/%v", DeploymentDirectoryFlag, v.Name()))
			checkError(err)
			defer f.Close()

			fileObjects, err := yamlReadAndSplit(f)
			checkError(err)

			for _, obj := range fileObjects {
				obj, objType, err := k8s.ParseManifest(obj)
				checkError(err)

				logrus.Info("Deploying: ", objType)

				err = k8sC.Create(ctx, obj.(ctrlclient.Object))
				checkError(err)
			}
		}
		logrus.Info("Directory Deployments Injected")

	}

	// err = docker.BuildAndPushImage(ctx, dockerfileString)
	// checkError(err)

	// // Create the Bootstrap Flux Resources
	// err = k8sC.Create(ctx, &ocirepo)
	// checkError(err)
	// err = k8sC.Create(ctx, &kustomizationOCI)
	// checkError(err)

	// // Wait for the flux
	// flux.WaitForKustomization(ctx, k8sC, kustomizationOCI.ObjectMeta)
}

func checkError(err error) {
	if err != nil {
		logrus.Fatal(err)
		panic(err)
	}
}
