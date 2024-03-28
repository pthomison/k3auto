package docker

import (
	"context"
	"path/filepath"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/pkg/archive"
)

func BuildImage(ctx context.Context, buildContextLocation string, dockerfile string, tags []string) error {
	apiClient, err := NewClient()
	if err != nil {
		return err
	}
	defer apiClient.Close()

	buildContextPath, err := filepath.Abs(buildContextLocation)
	if err != nil {
		return err
	}

	tarOpts := &archive.TarOptions{
		ExcludePatterns: []string{},
		IncludeFiles:    []string{"."},
		Compression:     archive.Uncompressed,
		NoLchown:        true,
	}
	tarArchive, err := archive.TarWithOptions(buildContextPath, tarOpts)
	if err != nil {
		return err
	}

	resp, err := apiClient.ImageBuild(ctx, tarArchive, types.ImageBuildOptions{
		Tags:       tags,
		Dockerfile: dockerfile,
	})
	if err != nil {
		return err
	}

	_ = resp

	return nil
}
