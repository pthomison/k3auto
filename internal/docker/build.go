package docker

import (
	"archive/tar"
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"io/fs"

	"github.com/docker/docker/api/types"
	"github.com/spf13/afero"
)

func BuildImage(ctx context.Context, buildContextLocation string, dockerfile string, tags []string, filesystem afero.Fs) (string, error) {
	apiClient, err := NewClient()
	if err != nil {
		return "", err
	}
	defer apiClient.Close()

	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)

	dockerFileLocation := fmt.Sprintf("%v/%v", buildContextLocation, "Dockerfile")

	hdr := &tar.Header{
		Name: dockerFileLocation,
		Mode: int64(0644),
		Size: int64(len([]byte(dockerfile))),
	}
	if err := tw.WriteHeader(hdr); err != nil {
		return "", err
	}

	_, err = tw.Write([]byte(dockerfile))
	if err != nil {
		return "", err
	}

	err = afero.Walk(filesystem, buildContextLocation, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		hdr := &tar.Header{
			Name: path,
			Mode: int64(info.Mode()),
			Size: info.Size(),
		}
		if err := tw.WriteHeader(hdr); err != nil {
			return err
		}

		f, err := filesystem.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		if _, err := io.Copy(tw, f); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return "", err
	}

	err = tw.Close()
	if err != nil {
		return "", err
	}

	resp, err := apiClient.ImageBuild(ctx, &buf, types.ImageBuildOptions{
		Tags:       tags,
		Dockerfile: dockerFileLocation,
	})
	if err != nil {
		return "", err
	}

	io.ReadAll(resp.Body)
	resp.Body.Close()

	hasher := sha256.New()
	hasher.Write([]byte(buf.Bytes()))
	hash := fmt.Sprintf("%x", hasher.Sum(nil))

	return hash, nil
}
