package cmd

import (
	"context"
	"errors"
	"io/fs"
	"path/filepath"
	"strings"

	defaults "github.com/pthomison/k3auto/default"
	"github.com/pthomison/k3auto/internal/flux"
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
		checkError(err)

		p := filepath.Join(exportDirectory, path)
		cp := filepath.Clean(strings.Replace(p, defaults.DefaultDeploymentsFolder, "", -1))

		if !info.IsDir() {
			b, err := afero.ReadFile(embedfs, path)
			checkError(err)

			if exists, _ := afero.Exists(osfs, cp); exists {
				err = osfs.Remove(cp)
				checkError(err)
			}

			err = afero.WriteFile(osfs, cp, b, info.Mode())
			checkError(err)
		} else {
			err := osfs.Mkdir(cp, 0755)
			if !errors.Is(err, afero.ErrFileExists) {
				checkError(err)
			}
		}

		return nil
	})

	manifests, err := flux.GenerateManifests("v2.3.0")
	checkError(err)

	p := filepath.Join(exportDirectory, manifests.Path)

	osfs.MkdirAll(filepath.Dir(p), 0755)

	err = afero.WriteFile(osfs, p, []byte(manifests.Content), 0644)
	checkError(err)
}
