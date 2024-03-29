package hack

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/pthomison/k3auto/internal/docker"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

const (
	dockerfile = "k3auto.Dockerfile"
	imageRef   = "k3auto-hack:latest"
	contextDir = ".."
)

// func TestInterfaces(t *testing.T) {
// 	ifaces, err := net.Interfaces()
// 	if err != nil {
// 		fmt.Print(fmt.Errorf("localAddresses: %+v\n", err.Error()))
// 		return
// 	}

// 	for _, i := range ifaces {
// 		addrs, err := i.Addrs()

// 		for _, a := range addrs {
// 			i, n, err := net.ParseCIDR(a.String())
// 			assert.Nil(t, err)

// 			if i.To4() != nil {
// 				spew.Dump(i, n)
// 			}

// 		}

// 		if i.Name == "en0" {
// 			fmt.Println(i.Name)
// 			// spew.Dump(addrs)

// 			for _, a := range addrs {
// 				i, n, err := net.ParseCIDR(a.String())
// 				assert.Nil(t, err)

// 				if i.To4() != nil {
// 					spew.Dump(i, n)
// 				}

// 			}
// 		}

// 		if err != nil {
// 			fmt.Print(fmt.Errorf("localAddresses: %+v\n", err.Error()))
// 			continue
// 		}
// 		for _, a := range addrs {
// 			switch v := a.(type) {
// 			case *net.IPAddr:
// 				fmt.Printf("%v : %s (%s)\n", i.Name, v, v.IP.DefaultMask())
// 			}

// 		}
// 	}
// }

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

	// image, err := docker.GetImageByName(ctx, imageRef)
	// assert.Nil(t, err)

	// spew.Dump(image)

	// hashRef := strings.ReplaceAll(image.RepoDigests[0], "@sha256:", ":")

	// err = docker.TagImage(ctx, image.RepoDigests[0], hashRef)
	// assert.Nil(t, err)

	// // spew.Dump(image)
	// // spew.Dump(imageRef)

	// err = docker.PushImage(ctx, hashRef, "127.0.0.1:8888")
	// assert.Nil(t, err)

}
