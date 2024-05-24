package cmd

import (
	"context"
	"errors"
	"io/fs"
	"path/filepath"

	"github.com/davecgh/go-spew/spew"
	defaults "github.com/pthomison/k3auto/default"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var (
	exportDirectory string
)

var ExportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export embeded resources",
	Run:   k3AutoExport,
}

func init() {
	ExportCmd.PersistentFlags().StringVarP(&exportDirectory, "export-directory", "e", ".", "Export Directory")

}

func k3AutoExport(cmd *cobra.Command, args []string) {
	ctx := cmd.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	osfs := afero.OsFs{}
	embedfs := afero.FromIOFS{FS: defaults.DefaultDeployments}

	afero.Walk(embedfs, ".", func(path string, info fs.FileInfo, err error) error {
		spew.Dump(path)
		checkError(err)

		p := filepath.Join(exportDirectory, path)
		spew.Dump(p)

		if !info.IsDir() {
			b, err2 := afero.ReadFile(embedfs, path)
			checkError(err2)

			err2 = afero.WriteFile(osfs, p, b, info.Mode())
			checkError(err2)
		} else {
			err2 := osfs.Mkdir(p, 0755)
			if !errors.Is(err2, afero.ErrFileExists) {
				checkError(err2)
			}
		}

		return nil
	})
}
