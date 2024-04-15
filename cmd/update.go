package cmd

import (
	"context"

	"github.com/pthomison/k3auto/pkg/k3auto"
	"github.com/spf13/cobra"
)

var UpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Reinject deployments",
	Run:   k3AutoUpdate,
}

func k3AutoUpdate(cmd *cobra.Command, args []string) {
	ctx := cmd.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	checkError(k3auto.Update(ctx, k3aConfig))
}
