package cmd

import (
	_ "embed"

	"github.com/pthomison/k3auto/pkg/k3auto"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var K3AutoCmd = &cobra.Command{
	Use:   "k3auto",
	Short: "k3auto is a local kubernetes cluster orchestrator powered by k3d and flux",
}

var (
	k3aConfig k3auto.Config = k3auto.Config{
		DeploymentFilesystem: afero.NewOsFs(),
	}
)

func init() {
	K3AutoCmd.PersistentFlags().StringVarP(&k3aConfig.ClusterConfigFile, "cluster-config", "c", "", "Override Cluster Config File")
	K3AutoCmd.PersistentFlags().StringVarP(&k3aConfig.SecretFile, "secret-config", "s", "", "Inject Secrets To the Cluster on Creation")
	K3AutoCmd.PersistentFlags().StringVarP(&k3aConfig.DeploymentDirectory, "deployment-directory", "d", "", "Deployment Directory")
	K3AutoCmd.PersistentFlags().StringVarP(&k3aConfig.BootstrapDirectory, "bootstrap-directory", "b", "/", "Folder Within The Deployment Directory To Bootstrap From")
	K3AutoCmd.PersistentFlags().BoolVarP(&k3aConfig.Minimal, "minimal", "m", false, "Only deploy the k3d cluster & flux controllers")

	K3AutoCmd.AddCommand(VersionCmd)
	K3AutoCmd.AddCommand(CreateCmd)
	K3AutoCmd.AddCommand(DeleteCmd)
	K3AutoCmd.AddCommand(UpdateCmd)
	// K3AutoCmd.AddCommand(ForwardCmd)
}

func checkError(err error) {
	if err != nil {
		logrus.Fatal(err)
		panic(err)
	}
}
