package main

import (
	"context"
	"log"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/fluxcd/flux2/v2/pkg/manifestgen/install"
	k3cluster "github.com/k3d-io/k3d/v5/pkg/client"
	"github.com/k3d-io/k3d/v5/pkg/config"
	"github.com/k3d-io/k3d/v5/pkg/runtimes"
	"github.com/sirupsen/logrus"
	"sigs.k8s.io/controller-runtime/pkg/client"
	ctlrconfig "sigs.k8s.io/controller-runtime/pkg/client/config"
)

var (
	rt = runtimes.Docker
)

func main() {
	ctx := context.TODO()

	err := config.ProcessSimpleConfig(clusterSimpleCfg)
	if err != nil {
		logrus.Fatal(err)
	}

	clusterConfig, err := config.TransformSimpleToClusterConfig(ctx, rt, *clusterSimpleCfg)
	if err != nil {
		logrus.Fatal(err)
	}

	if err := config.ValidateClusterConfig(ctx, runtimes.SelectedRuntime, *clusterConfig); err != nil {
		logrus.Fatal(err)
	}

	err = k3cluster.ClusterRun(ctx, rt, clusterConfig)
	if err != nil {
		logrus.Fatal(err)
	}

	if _, err := k3cluster.KubeconfigGetWrite(ctx, runtimes.SelectedRuntime,
		&clusterConfig.Cluster,
		"",
		&k3cluster.WriteKubeConfigOptions{
			UpdateExisting:       true,
			OverwriteExisting:    true,
			UpdateCurrentContext: clusterSimpleCfg.Options.KubeconfigOptions.SwitchCurrentContext,
		},
	); err != nil {
		logrus.Fatal(err)
	}

	time.Sleep(15 * time.Second)

	kcfg, err := ctlrconfig.GetConfig()
	if err != nil {
		logrus.Fatal(err)
	}

	logrus.Info("Get kubernetes config.")
	k8s, err := client.New(kcfg, client.Options{})
	if err != nil {
		log.Fatalf("Error creating kubernetes client: %s", err)
	}

	ready := false
	for !ready {
		ready, err = ArePodsReadyInCluster(ctx, k8s)
		if err != nil {
			logrus.Fatal(err)
		}
		logrus.Info("Waiting on cluster")
		time.Sleep(10 * time.Second)
	}

	spew.Dump(install.Generate(install.MakeDefaultOptions(), ""))

}
