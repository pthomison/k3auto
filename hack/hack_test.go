package hack

import (
	"context"
	"fmt"
	"net"
	"os"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/pthomison/k3auto/internal/docker"
	"github.com/pthomison/k3auto/internal/flux"
	"github.com/pthomison/k3auto/internal/k8s"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

const (
	dockerfile = "k3auto.Dockerfile"
	imageRef   = "k3auto-hack:latest"
	contextDir = ".."
)

func TestImageLookup(t *testing.T) {
	return
	ctx := context.TODO()

	apiClient, err := docker.NewClient()
	assert.Nil(t, err)
	defer apiClient.Close()

	err = os.WriteFile(fmt.Sprintf("%v/%v", contextDir, dockerfile), []byte(docker.DumpDockerfile), 0644)
	assert.Nil(t, err)
	defer os.Remove(fmt.Sprintf("%v/%v", contextDir, dockerfile))

	err = docker.BuildImage(ctx, contextDir, dockerfile, []string{imageRef}, afero.NewOsFs())
	spew.Dump(err)
	assert.Nil(t, err)

	err = docker.PushImage(ctx, imageRef, "127.0.0.1:8888")
	assert.Nil(t, err)
}

func TestKubectlApply(t *testing.T) {
	return
	ctx := context.TODO()

	fluxManifests, err := flux.GenerateManifests()
	assert.Nil(t, err)

	k8sC, err := k8s.NewClient()
	assert.Nil(t, err)

	k8s.CreateManifests(ctx, k8sC, fluxManifests.Content)

}

func TestIpLookup(t *testing.T) {
	// get list of available addresses
	addr, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, addr := range addr {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			// check if IPv4 or IPv6 is not nil
			if ipnet.IP.To4() != nil {
				// print available addresses
				fmt.Println(ipnet.IP.String())
			}
		}
	}
}
