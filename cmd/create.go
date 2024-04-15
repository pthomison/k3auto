package cmd

import (
	"context"

	"github.com/pthomison/k3auto/pkg/k3auto"
	"github.com/spf13/cobra"
	ctrl "sigs.k8s.io/controller-runtime"
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

func k3AutoCreate(cmd *cobra.Command, args []string) {
	ctx := cmd.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	_, err := k3auto.Create(ctx, k3aConfig)
	checkError(err)

}
