package cmd

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/pthomison/k3auto/pkg/k3auto"
)

var DeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete an existing cluster",
	Run:   k3AutoDelete,
}

func k3AutoDelete(cmd *cobra.Command, args []string) {
	ctx := cmd.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	checkError(k3auto.Delete(ctx, k3aConfig))
}
