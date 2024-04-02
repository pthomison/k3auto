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

// func PushImage(ctx context.Context, imageRef string, remoteRegistry string) error {
// 	apiClient, err := NewClient()
// 	if err != nil {
// 		return err
// 	}

// 	remoteRef := fmt.Sprintf("%v/%v", remoteRegistry, imageRef)

// 	err = apiClient.ImageTag(ctx, imageRef, remoteRef)
// 	if err != nil {
// 		return err
// 	}

// 	// ref, err := reference.ParseNormalizedNamed(remoteRef)
// 	// if err != nil {
// 	// 	return err
// 	// }

// 	// // Resolve the Repository name from fqn to RepositoryInfo
// 	// repoInfo, err := registry.ParseRepositoryInfo(ref)
// 	// if err != nil {
// 	// 	return err
// 	// }

// 	// Resolve the Auth config relevant for this server
// 	// authConfig := registry.AuthConfig{}
// 	// encodedAuth, err := registrytypes.EncodeAuthConfig(authConfig)
// 	// if err != nil {
// 	// 	return err
// 	// }
// 	// requestPrivilege := command.RegistryAuthenticationPrivilegedFunc(dockerCli, repoInfo.Index, "push")
// 	// options := types.ImagePushOptions{
// 	// 	All:           true,
// 	// 	RegistryAuth:  "",
// 	// 	PrivilegeFunc: func() (string, error) { return "", nil },
// 	// }

// 	// output, err := apiClient.ImagePush(ctx, remoteRef, options)
// 	// if err != nil {
// 	// 	return err
// 	// }

// 	cli, err := command.NewDockerCli()
// 	if err != nil {
// 		return err
// 	}

// 	pushCmd := image.NewPushCommand(cli)

// 	pushCmd.SetArgs([]string{remoteRef})

// 	err = pushCmd.Execute()
// 	if err != nil {
// 		return err
// 	}

// 	// io.ReadAll(output)

// 	return nil
// }
