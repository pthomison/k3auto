package secrets

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func InjectSecrets(ctx context.Context, k8sC client.Client, conf SecretConfig) error {

	binnedSecrets := map[string][]Secret{}

	for _, rawSec := range conf.Secrets {
		// Inject Default Fields If Not Present
		if rawSec.SecretName == "" {
			rawSec.SecretName = conf.DefaultName
		}
		if rawSec.SecretNamespace == "" {
			rawSec.SecretNamespace = conf.DefaultNamespace
		}

		// Collect Secrets Into K8S Secrets
		identifier := fmt.Sprintf("%v-%v", rawSec.SecretName, rawSec.SecretNamespace)
		if binnedSecrets[identifier] == nil {
			binnedSecrets[identifier] = []Secret{rawSec}
		} else {
			binnedSecrets[identifier] = append(binnedSecrets[identifier], rawSec)
		}
	}

	for _, secs := range binnedSecrets {
		k8sSecret := corev1.Secret{
			ObjectMeta: v1.ObjectMeta{
				Name:      secs[0].SecretName,
				Namespace: secs[0].SecretNamespace,
			},
		}

		sd := map[string]string{}

		for _, s := range secs {
			val, err := ResolverMap[s.Type].Resolve(ctx, s.Args)
			if err != nil {
				return err
			}
			sd[s.SecretKey] = val
		}

		k8sSecret.StringData = sd

		err := k8sC.Create(ctx, &k8sSecret)
		if err != nil {
			return err
		}
	}

	return nil
}
