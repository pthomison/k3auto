package main

import (
	"context"
	"os/exec"
	"time"

	"github.com/fluxcd/flux2/v2/pkg/manifestgen/install"
	"github.com/k3d-io/k3d/v5/pkg/runtimes"
	"github.com/sirupsen/logrus"
)

var (
	rt = runtimes.Docker
)

func checkError(err error) {
	if err != nil {
		logrus.Fatal(err)
	}
}

func main() {
	ctx := context.TODO()

	err := DeployCluster(ctx, clusterSimpleCfg)
	checkError(err)

	time.Sleep(15 * time.Second)

	k8s, err := k8sClient()
	checkError(err)

	ready := false
	for !ready {
		ready, err = ArePodsReadyInCluster(ctx, k8s)
		checkError(err)

		logrus.Info("Waiting on cluster")
		time.Sleep(10 * time.Second)
	}

	genOps := install.MakeDefaultOptions()
	genOps.NetworkPolicy = false
	fManifest, err := install.Generate(genOps, "")
	checkError(err)

	fileLoc, err := fManifest.WriteFile(".tmp")
	checkError(err)

	cmd := exec.Command("kubectl", "apply", "-f", fileLoc)
	err = cmd.Run()
	checkError(err)

}
