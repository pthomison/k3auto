package secrets

import (
	"context"
	"os/exec"
)

type ExecResolver struct{}

func (r *ExecResolver) Resolve(ctx context.Context, args []string) (string, error) {
	out, err := exec.CommandContext(ctx, args[0], args[1:]...).Output()
	if err != nil {
		return "", err
	}

	return string(out), nil
}
