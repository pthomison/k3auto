package main

import (
	"fmt"
	"os"

	"github.com/pthomison/k3auto/cmd"
)

func main() {
	if err := cmd.K3AutoCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
