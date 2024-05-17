package docker

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
)

func GetContainerByName(ctx context.Context, name string) (types.Container, error) {
	apiClient, err := NewClient()
	if err != nil {
		return types.Container{}, err
	}
	defer apiClient.Close()

	containers, err := apiClient.ContainerList(ctx, container.ListOptions{
		All: true,
	})
	if err != nil {
		return types.Container{}, err
	}

	for _, c := range containers {
		for _, n := range c.Names {
			if n == name {
				return c, nil
			}
		}
	}
	return types.Container{}, fmt.Errorf("could not find the docker container named %v", name)
}
