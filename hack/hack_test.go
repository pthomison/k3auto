package hack

import (
	"os"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/spf13/viper"

	kustomizev1 "github.com/fluxcd/kustomize-controller/api/v1"
	appsv1 "k8s.io/api/apps/v1"

	k3dconfig "github.com/k3d-io/k3d/v5/pkg/config"

	"github.com/stretchr/testify/assert"
)

const (
	deploymentLocation    = "./test-deployment.yaml"
	kustomizationLocation = "./test-kustomization.yaml"
	k3dConfigLocation     = "./test-k3dconfig.yaml"
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

	spew.Dump(cfg)
}

func TestCheckType(t *testing.T) {

}
