package hack

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/davecgh/go-spew/spew"
	defaults "github.com/pthomison/k3auto/default"
	"github.com/pthomison/k3auto/internal/flux"
	"github.com/spf13/afero"
	"github.com/spf13/viper"

	kustomizev1 "github.com/fluxcd/kustomize-controller/api/v1"
	appsv1 "k8s.io/api/apps/v1"

	k3dconfig "github.com/k3d-io/k3d/v5/pkg/config"

	"github.com/stretchr/testify/assert"
)

const (
	deploymentLocation    = "test-deployment.yaml"
	kustomizationLocation = "test-kustomization.yaml"
	k3dConfigLocation     = "test-k3dconfig.yaml"
	crdConfigLocation     = "test-crd.yaml"
)

func TestDecodeDeployment(t *testing.T) {
	yb, err := os.ReadFile(deploymentLocation)
	assert.Nil(t, err)
	obj, _, err := deserialize(yb)
	assert.Nil(t, err)

	assert.IsTypef(t, obj, &appsv1.Deployment{}, "Decoded Object Is Not appsv1.Deployment")
}

func TestDecodeKustomization(t *testing.T) {
	yb, err := os.ReadFile(kustomizationLocation)
	assert.Nil(t, err)
	obj, _, err := deserialize(yb)
	assert.Nil(t, err)

	assert.IsTypef(t, obj, &kustomizev1.Kustomization{}, "Decoded Object Is Not appsv1.Deployment")
}

func TestDecodeCRD(t *testing.T) {

	// spew.Dump(yb)

	// cmd := kubectl.NewDefaultKubectlCommand()
	// cmd.SetArgs([]string{"describe"})
	// err = cmd.Execute()
	// assert.Nil(t, err)

	// var crd apiextensionsv1

	// obj, _, err := deserialize(yb)
	// assert.Nil(t, err)

	// assert.IsTypef(t, obj, &kustomizev1.Kustomization{}, "Decoded Object Is Not appsv1.Deployment")
}

func TestDecodeK3dConfig(t *testing.T) {
	config := viper.New()
	config.SetConfigFile(k3dConfigLocation)

	if err := config.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			t.Error(err)
		}
		t.Error(err)
	}

	cfg, err := k3dconfig.FromViper(config)
	if err != nil {
		t.Error(err)
	}

	assert.NotNil(t, cfg)
}

func TestEmbededK3dConfig(t *testing.T) {
	config := viper.New()
	config.SetFs(afero.FromIOFS{FS: defaults.K3dConfig})
	config.SetConfigFile(defaults.K3dConfigLocation)

	if err := config.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			t.Error(err)
		}
		t.Error(err)
	}

	cfg, err := k3dconfig.FromViper(config)
	if err != nil {
		t.Error(err)
	}

	assert.NotNil(t, cfg)
}

func TestParseTypes(t *testing.T) {
	spew.Dump()

	files, err := defaults.DefaultDeployments.ReadDir("deployments")
	assert.Nil(t, err)

	for _, v := range files {
		spew.Dump(v)
		f, err := defaults.DefaultDeployments.Open(fmt.Sprintf("deployments/%v", v.Name()))
		assert.Nil(t, err)
		defer f.Close()

		fb, err := io.ReadAll(f)
		assert.Nil(t, err)

		objs := bytes.Split(fb, []byte("---"))

		for _, obj := range objs {
			if len(obj) != 0 {
				obj, objType, err := deserialize(obj)
				assert.Nil(t, err)

				_ = obj
				_ = objType
				// spew.Dump(obj, objType)
			}
		}

	}
}

func TestFluxManifests(t *testing.T) {
	return
	spew.Dump(flux.GenerateManifests())
}
