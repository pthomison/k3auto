package secrets

import (
	"context"
	"os"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/pthomison/k3auto/internal/k8s"
	"github.com/stretchr/testify/assert"
)

func TestConfigLoad(t *testing.T) {
	// return
	f, err := os.Open("config.test.yaml")
	assert.Nil(t, err)
	conf, err := LoadConfigFile(f)
	assert.Nil(t, err)

	k8sC, err := k8s.NewClient()
	assert.Nil(t, err)

	err = InjectSecrets(context.TODO(), k8sC, conf)
	assert.Nil(t, err)

	// spew.Dump(conf)
}

func TestExec(t *testing.T) {
	return
	er := &ExecResolver{}

	val, err := er.Resolve(context.TODO(), []string{"/bin/bash", "-c", "echo -n hello world"})
	assert.Nil(t, err)

	spew.Dump(val)
}
