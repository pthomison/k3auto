package main

import (
	"context"
	"os/exec"
	"time"

	"github.com/fluxcd/flux2/v2/pkg/manifestgen/install"
	"github.com/k3d-io/k3d/v5/pkg/runtimes"
	"github.com/sirupsen/logrus"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

var (
	rt = runtimes.Docker
)

func init() {
	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&zap.Options{})))
}

func main() {
	ctx := context.TODO()

	err := DeployCluster(ctx, clusterSimpleCfg)
	checkError(err)

	time.Sleep(10 * time.Second)

	k8s, err := k8sClient()
	checkError(err)

	err = WaitForPodsReadInCluster(ctx, k8s)
	checkError(err)

	genOps := install.MakeDefaultOptions()
	genOps.NetworkPolicy = false
	fManifest, err := install.Generate(genOps, "")
	checkError(err)

	fileLoc, err := fManifest.WriteFile(".tmp")
	checkError(err)

	cmd := exec.Command("kubectl", "apply", "-f", fileLoc)
	err = cmd.Run()
	checkError(err)

	err = WaitForPodsReadInCluster(ctx, k8s)
	checkError(err)

	err = k8s.Create(ctx, &secret)
	checkError(err)

	err = k8s.Create(ctx, &gitrepo)
	checkError(err)

	err = k8s.Create(ctx, &kustomization)
	checkError(err)
}

func checkError(err error) {
	if err != nil {
		logrus.Fatal(err)
	}
}
