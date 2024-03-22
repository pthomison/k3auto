package k8s

import (
	"bytes"
	"context"

	"github.com/davecgh/go-spew/spew"
	ctrl "sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	YamlDocSeperator = "---"
)

func Apply(ctx context.Context, k8sC ctrl.Client, b []byte) error {

	objs := bytes.Split(b, []byte(YamlDocSeperator))

	for _, obj := range objs {
		if len(obj) != 0 {
			o, t, err := ParseManifest(obj)
			if err != nil {
				return err
			}

			spew.Dump(o, t)

			err = k8sC.Create(ctx, o.(ctrl.Object))
			if err != nil {
				return err
			}
		}
	}

	return nil
}
