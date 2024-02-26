package main

import (
	"context"
	_ "embed"
	"os"
	"os/exec"

	"github.com/sirupsen/logrus"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/fluxcd/flux2/v2/pkg/manifestgen/install"

	"github.com/k3d-io/k3d/v5/pkg/runtimes"
)

var (
	rt = runtimes.Docker

	//go:embed Dockerfile
	DockerfileString string
)

func main() {
	ctx := context.TODO()

	// Deploy the cluster defined in cluster.go
	err := DeployCluster(ctx, clusterSimpleCfg)
	checkError(err)

	// Generate a k8s client from standard kubeconfig
	k8s, err := k8sClient()
	checkError(err)

	// Wait for the base cluster deployments to be ready
	WaitForDeployment(ctx, k8s, v1.ObjectMeta{
		Name:      "coredns",
		Namespace: "kube-system",
	})

	// Generate Flux Controller Manifests
	genOps := install.MakeDefaultOptions()
	genOps.NetworkPolicy = false
	fManifest, err := install.Generate(genOps, "")
	checkError(err)

	// Write Controller Manifests to tmp folder
	fileLoc, err := fManifest.WriteFile(os.TempDir())
	checkError(err)
	defer os.Remove(fileLoc)

	// Apply Controller Manifests
	// TODO: Figure out a way to do this w/o exec & kubectl
	cmd := exec.Command("kubectl", "apply", "-f", fileLoc)
	err = cmd.Run()
	checkError(err)

	err = BuildAndPushImage(ctx)
	checkError(err)

	// Create the Bootstrap Flux Resources
	err = k8s.Create(ctx, &ocirepo)
	checkError(err)
	err = k8s.Create(ctx, &kustomizationOCI)
	checkError(err)

	// Wait for the flux
	WaitForKustomization(ctx, k8s, kustomizationOCI.ObjectMeta)

}

func checkError(err error) {
	if err != nil {
		logrus.Fatal(err)
	}
}
