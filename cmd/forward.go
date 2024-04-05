package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var ForwardCmd = &cobra.Command{
	Use:   "forward",
	Short: "Forwards Ports To the Environment",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Implementation TBD")
	},
}
