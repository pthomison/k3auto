package hack

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/pthomison/k3auto/internal/flux"
	"github.com/stretchr/testify/assert"
)

func TestManifestGeneration(t *testing.T) {
	manifests, err := flux.GenerateManifests("v2.3.0")
	assert.Nil(t, err)

	spew.Dump(manifests)

}
