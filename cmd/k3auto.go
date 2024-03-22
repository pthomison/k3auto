package cmd

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"io"
	"os"
	"path"
	"time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	k3dv1alpha5 "github.com/k3d-io/k3d/v5/pkg/config/v1alpha5"
	"github.com/k3d-io/k3d/v5/pkg/runtimes"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	defaults "github.com/pthomison/k3auto/default"
	"github.com/pthomison/k3auto/internal/flux"
	"github.com/pthomison/k3auto/internal/k3d"
	"github.com/pthomison/k3auto/internal/k8s"

	kubectl "k8s.io/kubectl/pkg/cmd"
)

var K3AutoCmd = &cobra.Command{
	Use: "k3auto",
	// Short: "Hugo is a very fast static site generator",
	// Long: `A Fast and Flexible Static Site Generator built with
	// 			  love by spf13 and friends in Go.
	// 			  Complete documentation is available at https://gohugo.io/documentation/`,
	Run: k3AutoRun,
}

const ()

var (
	ClusterConfigFileFlag   string
	DeploymentDirectoryFlag string
	MinimalFlag             bool
)

func init() {
	K3AutoCmd.PersistentFlags().StringVarP(&ClusterConfigFileFlag, "cluster-config", "c", "", "Override Cluster Config File")
	K3AutoCmd.PersistentFlags().StringVarP(&DeploymentDirectoryFlag, "deployment-directory", "d", "", "Deployment Directory")
	K3AutoCmd.PersistentFlags().BoolVarP(&MinimalFlag, "minimal", "m", false, "Only deploy the k3d cluster")

	opts := zap.Options{
		Development: true,
	}
	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

}

func k3AutoRun(cmd *cobra.Command, args []string) {

	ctx := cmd.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	rt := runtimes.Docker

	var clusterConfig *k3dv1alpha5.SimpleConfig
	var err error

	if ClusterConfigFileFlag != "" {
		clusterConfig, err = k3d.ParseConfigFile(ClusterConfigFileFlag, nil)
		checkError(err)
	} else {
		clusterConfig, err = k3d.ParseConfigFile(defaults.K3dConfigLocation, &defaults.K3dConfig)
		checkError(err)
	}

	// Deploy the cluster defined in cluster.go
	err = k3d.DeployCluster(ctx, clusterConfig, rt)
	checkError(err)

	// Generate a k8s client from standard kubeconfig
	k8sC, err := k8s.NewClient()
	checkError(err)

	// Wait for the base cluster deployments to be ready
	k8s.WaitForDeployment(ctx, k8sC, v1.ObjectMeta{
		Name:      "coredns",
		Namespace: "kube-system",
	})

	// Generate Flux Controller Manifests
	fluxManifests, err := flux.GenerateManifests()
	checkError(err)

	tmpDirLoc, err := os.MkdirTemp("", "k3auto-")
	checkError(err)
	defer os.RemoveAll(tmpDirLoc)

	fluxManifestsPath := path.Join(tmpDirLoc, "flux-manifests.yaml")
	os.WriteFile(fluxManifestsPath, []byte(fluxManifests.Content), 0644)

	kubectlCmd := kubectl.NewDefaultKubectlCommand()
	kubectlCmd.SetArgs([]string{"apply", "-f", fluxManifestsPath})
	err = kubectlCmd.Execute()
	checkError(err)
	time.Sleep(1 * time.Second)

	deploymentFiles, err := defaults.DefaultDeployments.ReadDir("deployments")
	checkError(err)

	for _, v := range deploymentFiles {
		f, err := defaults.DefaultDeployments.Open(fmt.Sprintf("deployments/%v", v.Name()))
		checkError(err)
		defer f.Close()

		fb, err := io.ReadAll(f)
		checkError(err)

		objs := bytes.Split(fb, []byte("---"))

		for _, obj := range objs {
			if len(obj) != 0 {
				obj, objType, err := k8s.ParseManifest(obj)
				checkError(err)

				_ = obj
				_ = objType
				logrus.Info("Deploying: ", objType)

				err = k8sC.Create(ctx, obj.(client.Object))
				checkError(err)
			}
		}
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
