package main

import (
	"context"
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
	k3dclient "github.com/k3d-io/k3d/v5/pkg/client"
	"github.com/k3d-io/k3d/v5/pkg/runtimes"
	"github.com/k3d-io/k3d/v5/pkg/types"
	"github.com/pthomison/k3auto/cmd"
	"github.com/pthomison/k3auto/internal/k8s"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

var (
	rt = runtimes.Docker
)

func init() {
	opts := zap.Options{
		Development: true,
	}
	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))
}

func SetupEnvironment(ctx context.Context) (func(ctx context.Context) error, error) {
	// go cmd.CreateCmd.ExecuteContext(ctx)

	cmd.K3AutoCmd.SetArgs([]string{"create", "-d", "./e2e_deployments", "-b", "./e2e_deployments/bootstrap/", "-s", "./e2e_secrets.yaml"})
	go cmd.K3AutoCmd.ExecuteContext(ctx)

	cleanupFn := func(ctx context.Context) error {
		cluster, err := k3dclient.ClusterGet(ctx, rt, &types.Cluster{
			Name: "k3auto",
		})
		if err != nil {
			return err
		}

		err = k3dclient.ClusterDelete(ctx, rt, cluster, types.ClusterDeleteOpts{})
		if err != nil {
			return err
		}

		return nil
	}

	for {
		_, k8sC, err := k8s.NewClient()
		if err == nil {
			p := corev1.PodList{}
			err = k8sC.List(ctx, &p)

			if err == nil {
				break
			}
		}

		time.Sleep(1 * time.Second)
	}
	return cleanupFn, nil
}

func TestEndToEnd(t *testing.T) {
	logrus.Info("Starting End to End Test")

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(300*time.Second))
	defer cancel()
	go func() {
		select {
		case <-ctx.Done():
			switch ctx.Err() {
			case context.DeadlineExceeded:
				assert.FailNow(t, "Environment Timed-Out While Deploying")
				logrus.Fatal("fail")
			case context.Canceled:
				_ = ""
			}
		}
	}()

	cleanupFn, err := SetupEnvironment(ctx)
	defer cleanupFn(ctx)
	assert.Nil(t, err)

	deploymentList := appsv1.DeploymentList{}
	_, k8sC, err := k8s.NewClient()
	assert.Nil(t, err)

	err = k8sC.List(ctx, &deploymentList)
	assert.Nil(t, err)

	for _, dep := range deploymentList.Items {
		spew.Dump(dep.Name, dep.Namespace)
	}

	err = k8s.WaitForDeployment(ctx, k8sC, metav1.ObjectMeta{
		Name:      "coredns",
		Namespace: "kube-system",
	})
	assert.Nil(t, err)

	err = k8s.WaitForDeployment(ctx, k8sC, metav1.ObjectMeta{
		Name:      "helm-controller",
		Namespace: "flux-system",
	})
	assert.Nil(t, err)

	err = k8s.WaitForDeployment(ctx, k8sC, metav1.ObjectMeta{
		Name:      "kustomize-controller",
		Namespace: "flux-system",
	})
	assert.Nil(t, err)

	err = k8s.WaitForDeployment(ctx, k8sC, metav1.ObjectMeta{
		Name:      "metrics-server",
		Namespace: "metrics-server",
	})
	assert.Nil(t, err)

	err = k8s.WaitForDeployment(ctx, k8sC, metav1.ObjectMeta{
		Name:      "kube-state-metrics",
		Namespace: "kube-state-metrics",
	})
	assert.Nil(t, err)

}
