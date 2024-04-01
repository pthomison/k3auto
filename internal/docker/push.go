package docker

import (
	"context"
	"fmt"

	"github.com/containers/image/v5/copy"
	"github.com/containers/image/v5/signature"
	"github.com/containers/image/v5/transports/alltransports"
	containertypes "github.com/containers/image/v5/types"
)

func PushImage(ctx context.Context, imageRef string, remoteRegistry string) error {
	srcImage, err := alltransports.ParseImageName(fmt.Sprintf("docker-daemon:%v", imageRef))
	if err != nil {
		return err
	}

	destImage, err := alltransports.ParseImageName(fmt.Sprintf("docker://%v/%v", remoteRegistry, imageRef))
	if err != nil {
		return err
	}

	policy := &signature.Policy{Default: []signature.PolicyRequirement{signature.NewPRInsecureAcceptAnything()}}
	dpc, err := signature.NewPolicyContext(policy)
	if err != nil {
		return err
	}

	_, err = copy.Image(ctx, dpc, destImage, srcImage, &copy.Options{
		DestinationCtx: &containertypes.SystemContext{
			DockerDaemonInsecureSkipTLSVerify: true,
			OCIInsecureSkipTLSVerify:          true,
			DockerInsecureSkipTLSVerify:       containertypes.NewOptionalBool(true),
		},
		SourceCtx: &containertypes.SystemContext{
			DockerDaemonInsecureSkipTLSVerify: true,
			OCIInsecureSkipTLSVerify:          true,
			DockerInsecureSkipTLSVerify:       containertypes.NewOptionalBool(true),
		},
	})
	if err != nil {
		return err
	}

	return nil
}
