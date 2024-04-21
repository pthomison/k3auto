package flux

import (
	"context"
	"fmt"
	"time"

	"github.com/fluxcd/cli-utils/pkg/kstatus/status"
	"github.com/fluxcd/pkg/runtime/patch"
	"sigs.k8s.io/controller-runtime/pkg/client"

	kustomizev1 "github.com/fluxcd/kustomize-controller/api/v1"
	sourcev1beta2 "github.com/fluxcd/source-controller/api/v1beta2"

	apimachinerytypes "k8s.io/apimachinery/pkg/types"

	"k8s.io/client-go/util/retry"
)

func ReconcileKustomization(ctx context.Context, k8sC client.Client, name apimachinerytypes.NamespacedName, timeout time.Duration) (chan error, error) {

	requestTime := time.Now()

	err := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		k := kustomizev1.Kustomization{}
		err := k8sC.Get(ctx, name, &k)
		if err != nil {
			return err
		}

		if k.Annotations == nil {
			k.Annotations = make(map[string]string)
		}

		k.Annotations["reconcile.fluxcd.io/requestedAt"] = fmt.Sprintf("%v", requestTime.Unix())

		err = k8sC.Update(ctx, &k)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	timeoutCtx, cancelFunc := context.WithTimeout(ctx, timeout)
	completeChan := make(chan error)

	go func(ctx context.Context, ctxCancel context.CancelFunc, k8sC client.Client, name apimachinerytypes.NamespacedName, requestTime time.Time) {
		time.Sleep(time.Second * 1)

		for {
			k := kustomizev1.Kustomization{}
			err := k8sC.Get(ctx, name, &k)
			if err != nil {
				completeChan <- err
				close(completeChan)
				return
			}

			u, err := patch.ToUnstructured(&k)
			if err != nil {
				completeChan <- err
				close(completeChan)
				return
			}

			res, err := status.Compute(u)
			if err != nil {
				completeChan <- err
				close(completeChan)
				return
			}

			if res.Status == status.CurrentStatus {
				completeChan <- nil
				close(completeChan)
				cancelFunc()
				return
			}
		}

	}(timeoutCtx, cancelFunc, k8sC, name, requestTime)

	return completeChan, nil
}

func ReconcileOCIRepository(ctx context.Context, k8sC client.Client, name apimachinerytypes.NamespacedName, timeout time.Duration) (chan error, error) {

	requestTime := time.Now()

	err := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		repo := sourcev1beta2.OCIRepository{}
		err := k8sC.Get(ctx, name, &repo)
		if err != nil {
			return err
		}

		if repo.Annotations == nil {
			repo.Annotations = make(map[string]string)
		}

		repo.Annotations["reconcile.fluxcd.io/requestedAt"] = fmt.Sprintf("%v", requestTime.Unix())

		err = k8sC.Update(ctx, &repo)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	timeoutCtx, cancelFunc := context.WithTimeout(ctx, timeout)
	completeChan := make(chan error)

	go func(ctx context.Context, ctxCancel context.CancelFunc, k8sC client.Client, name apimachinerytypes.NamespacedName, requestTime time.Time) {
		for {
			time.Sleep(time.Second * 1)

			repo := sourcev1beta2.OCIRepository{}
			err := k8sC.Get(ctx, name, &repo)
			if err != nil {
				completeChan <- err
				close(completeChan)
				return
			}

			u, err := patch.ToUnstructured(&repo)
			if err != nil {
				completeChan <- err
				close(completeChan)
				return
			}

			res, err := status.Compute(u)
			if err != nil {
				completeChan <- err
				close(completeChan)
				return
			}

			if res.Status == status.CurrentStatus {
				completeChan <- nil
				close(completeChan)
				cancelFunc()
				return
			}
		}

	}(timeoutCtx, cancelFunc, k8sC, name, requestTime)

	return completeChan, nil
}
