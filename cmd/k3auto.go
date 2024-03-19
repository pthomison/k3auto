package cmd

import (
	"context"
	_ "embed"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/k3d-io/k3d/v5/pkg/runtimes"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/pthomison/k3auto/internal/k3d"
	"github.com/pthomison/k3auto/internal/k8s"
)

var K3AutoCmd = &cobra.Command{
	Use: "k3auto",
	// Short: "Hugo is a very fast static site generator",
	// Long: `A Fast and Flexible Static Site Generator built with
	// 			  love by spf13 and friends in Go.
	// 			  Complete documentation is available at https://gohugo.io/documentation/`,
	Run: k3AutoRun,
}

var (
	ClusterConfigFileFlag string
)

func init() {
	// cobra.OnInitialize(initConfig)

	K3AutoCmd.PersistentFlags().StringVar(&ClusterConfigFileFlag, "cluster-config", "./cluster-config.yaml", "")
	// rootCmd.PersistentFlags().StringP("author", "a", "YOUR NAME", "author name for copyright attribution")
	// rootCmd.PersistentFlags().StringVarP(&userLicense, "license", "l", "", "name of license for the project")
	// rootCmd.PersistentFlags().Bool("viper", true, "use Viper for configuration")
	// viper.BindPFlag("author", rootCmd.PersistentFlags().Lookup("author"))
	// viper.BindPFlag("useViper", rootCmd.PersistentFlags().Lookup("viper"))
	// viper.SetDefault("author", "NAME HERE <EMAIL ADDRESS>")
	// viper.SetDefault("license", "apache")

	// rootCmd.AddCommand(addCmd)
	// rootCmd.AddCommand(initCmd)
}

func k3AutoRun(cmd *cobra.Command, args []string) {
	ctx := cmd.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	rt := runtimes.Docker

	// Deploy the cluster defined in cluster.go
	err := k3d.DeployCluster(ctx, clusterSimpleCfg, rt)
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
	// genOps := install.MakeDefaultOptions()
	// genOps.NetworkPolicy = false
	// fManifest, err := install.Generate(genOps, "")
	// checkError(err)

	// // Write Controller Manifests to tmp folder
	// fileLoc, err := fManifest.WriteFile(os.TempDir())
	// checkError(err)
	// defer os.Remove(fileLoc)

	// // Apply Controller Manifests
	// // TODO: Figure out a way to do this w/o exec & kubectl
	// cmd := exec.Command("kubectl", "apply", "-f", fileLoc)
	// err = cmd.Run()
	// checkError(err)

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
	}
}
