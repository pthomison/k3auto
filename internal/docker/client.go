package docker

import dockerclient "github.com/docker/docker/client"

func NewClient() (*dockerclient.Client, error) {
	apiClient, err := dockerclient.NewClientWithOpts(dockerclient.FromEnv)
	if err != nil {
		return nil, err
	}

	return apiClient, err
	// defer apiClient.Close()
}
