package k8s

import (
	"time"

	kubectl "k8s.io/kubectl/pkg/cmd"
)

func KubectlApply(filepath string) error {
	kubectlCmd := kubectl.NewDefaultKubectlCommand()
	kubectlCmd.SetArgs([]string{"apply", "-f", filepath})
	err := kubectlCmd.Execute()
	if err != nil {
		return err
	}
	time.Sleep(1 * time.Second)
	return nil
}
