package main

import (
	"context"
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
	k3dclient "github.com/k3d-io/k3d/v5/pkg/client"
	"github.com/k3d-io/k3d/v5/pkg/types"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	appsv1 "k8s.io/api/apps/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	apimachinerytypes "k8s.io/apimachinery/pkg/types"

	ctrl "sigs.k8s.io/controller-runtime"
	logzap "sigs.k8s.io/controller-runtime/pkg/log/zap"
)

func init() {
	ctrl.SetLogger(logzap.New(logzap.UseFlagOptions(&logzap.Options{
		Development: true, // a sane default
		ZapOpts:     []zap.Option{zap.AddCaller()},
	})))
	// ctrl.SetLogger(logrus.New())
}

func SetupEnvironment(ctx context.Context) (func(ctx context.Context) error, error) {
	go main()

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

	time.Sleep(30 * time.Second)
	return cleanupFn, nil
}

func DeploymentReady(ctx context.Context, k8s client.Client, name string, namespace string) (bool, error) {
	dp := appsv1.Deployment{}
	err := k8s.Get(ctx, apimachinerytypes.NamespacedName{
		Name:      name,
		Namespace: namespace,
	}, &dp)
	if err != nil {
		return false, err
	}

	return (dp.Status.ReadyReplicas == *dp.Spec.Replicas), nil
}

func TestEndToEnd(t *testing.T) {
	logrus.Info("Starting End to End Test")

	ctx := context.TODO()
	cleanupFn, err := SetupEnvironment(ctx)
	defer cleanupFn(ctx)
	_ = cleanupFn
	assert.Nil(t, err)

	k8s, err := k8sClient()
	assert.Nil(t, err)

	deploymentList := appsv1.DeploymentList{}
	err = k8s.List(ctx, &deploymentList)
	assert.Nil(t, err)

	time.Sleep(10 * time.Second)

	for _, dep := range deploymentList.Items {
		spew.Dump(dep.Name, dep.Namespace)
	}

	hcReady, err := DeploymentReady(ctx, k8s, "helm-controller", "flux-system")
	assert.Nil(t, err)
	assert.True(t, hcReady)

	kcReady, err := DeploymentReady(ctx, k8s, "kustomize-controller", "flux-system")
	assert.Nil(t, err)
	assert.True(t, kcReady)

	ncReady, err := DeploymentReady(ctx, k8s, "notification-controller", "flux-system")
	assert.Nil(t, err)
	assert.True(t, ncReady)

	scReady, err := DeploymentReady(ctx, k8s, "source-controller", "flux-system")
	assert.Nil(t, err)
	assert.True(t, scReady)
}
