package hack

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/pthomison/k3auto/internal/docker"
	"github.com/stretchr/testify/assert"
)

const (
	deploymentLocation    = "test-deployment.yaml"
	kustomizationLocation = "test-kustomization.yaml"
	k3dConfigLocation     = "test-k3dconfig.yaml"
	crdConfigLocation     = "test-crd.yaml"
)

// func TestDecodeDeployment(t *testing.T) {
// 	yb, err := os.ReadFile(deploymentLocation)
// 	assert.Nil(t, err)
// 	obj, _, err := deserialize(yb)
// 	assert.Nil(t, err)

// 	assert.IsTypef(t, obj, &appsv1.Deployment{}, "Decoded Object Is Not appsv1.Deployment")
// }

// func TestDecodeKustomization(t *testing.T) {
// 	yb, err := os.ReadFile(kustomizationLocation)
// 	assert.Nil(t, err)
// 	obj, _, err := (yb)
// 	assert.Nil(t, err)

// 	assert.IsTypef(t, obj, &kustomizev1.Kustomization{}, "Decoded Object Is Not appsv1.Deployment")
// }

// func TestDecodeK3dConfig(t *testing.T) {
// 	config := viper.New()
// 	config.SetConfigFile(k3dConfigLocation)

// 	if err := config.ReadInConfig(); err != nil {
// 		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
// 			t.Error(err)
// 		}
// 		t.Error(err)
// 	}

// 	cfg, err := k3dconfig.FromViper(config)
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	assert.NotNil(t, cfg)
// }

// func TestEmbededK3dConfig(t *testing.T) {
// 	config := viper.New()
// 	config.SetFs(afero.FromIOFS{FS: defaults.K3dConfig})
// 	config.SetConfigFile(defaults.K3dConfigLocation)

// 	if err := config.ReadInConfig(); err != nil {
// 		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
// 			t.Error(err)
// 		}
// 		t.Error(err)
// 	}

// 	cfg, err := k3dconfig.FromViper(config)
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	assert.NotNil(t, cfg)
// }

// func TestParseTypes(t *testing.T) {
// 	spew.Dump()

// 	files, err := defaults.DefaultDeployments.ReadDir("deployments")
// 	assert.Nil(t, err)

// 	for _, v := range files {
// 		spew.Dump(v)
// 		f, err := defaults.DefaultDeployments.Open(fmt.Sprintf("deployments/%v", v.Name()))
// 		assert.Nil(t, err)
// 		defer f.Close()

// 		fb, err := io.ReadAll(f)
// 		assert.Nil(t, err)

// 		objs := bytes.Split(fb, []byte("---"))

// 		for _, obj := range objs {
// 			if len(obj) != 0 {
// 				obj, objType, err := deserialize(obj)
// 				assert.Nil(t, err)

// 				_ = obj
// 				_ = objType
// 				// spew.Dump(obj, objType)
// 			}
// 		}

// 	}
// }

func TestDocker(t *testing.T) {
	ctx := context.Background()

	apiClient, err := docker.NewClient()
	assert.Nil(t, err)
	defer apiClient.Close()

	// spew.Dump(apiClient.ContainerList(ctx, types.ContainerListOptions{
	// 	All: true,
	// }))

	err = os.WriteFile("../k3auto.Dockerfile", []byte(docker.DumpDockerfile), 0644)
	assert.Nil(t, err)
	defer os.Remove("../k3auto.Dockerfile")

	err = docker.BuildImage(ctx, "..", "k3auto.Dockerfile", []string{"k3auto-hack:latest"})
	spew.Dump(err)
	assert.Nil(t, err)

	// spew.Dump(docker.GetContainerByName(ctx, "/k3auto-registry"))

	err = docker.PushImage(ctx, "k3auto:latest", "127.0.0.1:8888")
	assert.Nil(t, err)

	// buildContext, err := docker.CreateTarStream(".", "k3auto.Dockerfile")
	// assert.Nil(t, err)

	// F := afero.NewOsFs()
	// for _, v := range []string{} {
	// 	F.W
	// }

	// resp, err := apiClient.ImageBuild(context.TODO(), buildContext, types.ImageBuildOptions{
	// 	Tags:       []string{"k3auto-fluxdir:latest"},
	// 	Dockerfile: "k3auto.Dockerfile",
	// })
	// assert.Nil(t, err)

	// spew.Dump(resp)
}

func TestImageLookup(t *testing.T) {
	ctx := context.TODO()

	image, err := docker.GetImageByName(ctx, "fedora:38")
	assert.Nil(t, err)

	// imageName := "fedora"
	// imageTag := image.RepoDigests[0]
	// imageRef := fmt.Sprintf("%v@%v", imageName, imageTag)

	// imageRef := image.RepoDigests[0]

	imageRef := strings.ReplaceAll(image.RepoDigests[0], "@sha256:", ":")

	err = docker.TagImage(ctx, image.RepoDigests[0], imageRef)
	assert.Nil(t, err)

	spew.Dump(image)
	spew.Dump(imageRef)

	err = docker.PushImage(ctx, imageRef, "127.0.0.1:8888")
	assert.Nil(t, err)

}
