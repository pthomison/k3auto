package cmd

import (
	_ "embed"

	"github.com/spf13/cobra"
)

var K3AutoCmd = &cobra.Command{
	Use:   "k3auto",
	Short: "k3auto is a local kubernetes cluster orchestrator powered by k3d and flux",
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

	K3AutoCmd.AddCommand(VersionCmd)
	K3AutoCmd.AddCommand(CreateCmd)
	K3AutoCmd.AddCommand(DeleteCmd)
	// K3AutoCmd.AddCommand(UpdateCmd)
}
