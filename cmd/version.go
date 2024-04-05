package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// ldflag vars
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints the version, commit, & build date",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Version: %v\nCommit: %v\nDate: %v\n", version, commit, date)
	},
}
