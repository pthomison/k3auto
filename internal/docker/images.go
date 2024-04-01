package docker

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
)

func GetImageByName(ctx context.Context, name string) (types.ImageSummary, error) {
	apiClient, err := NewClient()
	if err != nil {
		return types.ImageSummary{}, err
	}
	defer apiClient.Close()

	images, err := apiClient.ImageList(ctx, types.ImageListOptions{
		All: true,
	})
	if err != nil {
		return types.ImageSummary{}, err
	}

	for _, image := range images {
		for _, tag := range image.RepoTags {
			if tag == name {
				return image, nil
			}
		}
	}

	return types.ImageSummary{}, fmt.Errorf("could not find the docker image named %v", name)
}

func TagImage(ctx context.Context, src string, dest string) error {
	apiClient, err := NewClient()
	if err != nil {
		return err
	}
	defer apiClient.Close()

	err = apiClient.ImageTag(ctx, src, dest)
	if err != nil {
		return err
	}
	return nil
}
