package presenters

import (
	"errors"
	"fmt"
	"strings"

	"github.com/eleven-sh/cli/internal/config"
	"github.com/eleven-sh/cli/internal/exceptions"
	"github.com/eleven-sh/cli/internal/system"
	"github.com/eleven-sh/cli/internal/views"
	"github.com/eleven-sh/eleven/entities"
	"github.com/google/go-github/v43/github"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ViewableErrorBuilder interface {
	Build(error) *views.ViewableError
}

type ElevenViewableErrorBuilder struct{}

func NewElevenViewableErrorBuilder() ElevenViewableErrorBuilder {
	return ElevenViewableErrorBuilder{}
}

func (ElevenViewableErrorBuilder) Build(err error) (viewableError *views.ViewableError) {
	viewableError = &views.ViewableError{}

	if typedError, ok := err.(entities.ErrClusterNotExists); ok {
		viewableError.Title = "Cluster not found"
		viewableError.Message = fmt.Sprintf(
			"The cluster \"%s\" was not found.",
			config.ColorsBold(typedError.ClusterName),
		)

		return
	}

	if typedError, ok := err.(entities.ErrClusterAlreadyExists); ok {
		viewableError.Title = "Cluster already exists"
		viewableError.Message = fmt.Sprintf(
			"The cluster \"%s\" already exists.",
			config.ColorsBold(typedError.ClusterName),
		)

		return
	}

	if typedError, ok := err.(entities.ErrEnvNotExists); ok {
		viewableError.Title = "Sandbox not found"

		if typedError.ClusterName != entities.DefaultClusterName {
			viewableError.Message = fmt.Sprintf(
				"The sandbox \"%s\" was not found in the cluster \"%s\".",
				config.ColorsBold(typedError.EnvName),
				config.ColorsBold(typedError.ClusterName),
			)
			return
		}

		viewableError.Message = fmt.Sprintf(
			"The sandbox \"%s\" was not found.",
			config.ColorsBold(typedError.EnvName),
		)
		return
	}

	if errors.Is(err, exceptions.ErrUserNotLoggedIn) {
		viewableError.Title = "GitHub account not connected"
		viewableError.Message = fmt.Sprintf(
			"You must first connect your GitHub account using the command \"eleven login\".\n\n"+
				"Eleven requires the following permissions:\n\n"+
				"  - \"Public SSH keys\" and \"Repositories\" to let you access your repositories from your sandboxes\n\n"+
				"  - \"Personal user data\" to configure Git\n\n"+
				"All your data (including the OAuth access token) will only be stored locally (in \"%s\").",
			config.ColorsBold(system.UserConfigFilePath()),
		)

		return
	}

	if typedError, ok := err.(entities.ErrEnvRepositoryNotFound); ok {
		viewableError.Title = "Repository not found"
		viewableError.Message = fmt.Sprintf(
			"The repository \"%s\" was not found.\n\n"+
				"Please double check that this repository exists and that you can access it.",
			config.ColorsBold(typedError.RepoOwner+"/"+typedError.RepoName),
		)

		return
	}

	if typedError, ok := err.(entities.ErrEnvDuplicatedRepositories); ok {
		viewableError.Title = "Duplicated repositories"
		viewableError.Message = fmt.Sprintf(
			"The repository \"%s\" is set multiple times.",
			config.ColorsBold(typedError.RepoOwner+"/"+typedError.RepoName),
		)

		return
	}

	if typedError, ok := err.(entities.ErrEnvInvalidRuntime); ok {
		viewableError.Title = "Invalid runtime"
		viewableError.Message = fmt.Sprintf(
			"The runtime \"%s\" is not supported by Eleven.",
			config.ColorsBold(typedError.Runtime),
		)

		return
	}

	if typedError, ok := err.(entities.ErrEnvInvalidRuntimeVersion); ok {
		viewableError.Title = "Invalid runtime version"
		viewableError.Message = fmt.Sprintf(
			"The value \"%s\" is not a valid version for \"%s\".\n\n"+
				"Example(s) of valid version(s): %s",
			config.ColorsBold(typedError.RuntimeVersion),
			config.ColorsBold(typedError.Runtime),
			config.ColorsBold(strings.Join(typedError.RuntimeVersionExamples, ", ")),
		)

		return
	}

	if typedError, ok := err.(entities.ErrEnvDuplicatedRuntimes); ok {
		viewableError.Title = "Duplicated runtimes"
		viewableError.Message = fmt.Sprintf(
			"The runtime \"%s\" is set multiple times.",
			config.ColorsBold(typedError.Runtime),
		)

		return
	}

	if typedError, ok := err.(entities.ErrInitRemovingEnv); ok {
		viewableError.Title = "Invalid sandbox state"
		viewableError.Message = fmt.Sprintf(
			"The sandbox \"%s\" cannot be initialized because it's currently removing.\n\n"+
				"You must wait for the removing process to terminate.",
			config.ColorsBold(typedError.EnvName),
		)

		return
	}

	if typedError, ok := err.(entities.ErrEditRemovingEnv); ok {
		viewableError.Title = "Invalid sandbox state"
		viewableError.Message = fmt.Sprintf(
			"The sandbox \"%s\" cannot be edited because it's currently removing.\n\n"+
				"You must wait for the removing process to terminate.",
			config.ColorsBold(typedError.EnvName),
		)

		return
	}

	if typedError, ok := err.(entities.ErrEditCreatingEnv); ok {
		viewableError.Title = "Invalid sandbox state"
		viewableError.Message = fmt.Sprintf(
			"The sandbox \"%s\" cannot be edited because it's currently initializing.\n\n"+
				"You must wait for the initialization process to terminate.",
			config.ColorsBold(typedError.EnvName),
		)

		return
	}

	if typedError, ok := err.(entities.ErrServeCreatingEnv); ok {
		viewableError.Title = "Invalid sandbox state"
		viewableError.Message = fmt.Sprintf(
			"No port could be served in the sandbox \"%s\" because it's currently initializing.\n\n"+
				"You must wait for the initialization process to terminate.",
			config.ColorsBold(typedError.EnvName),
		)

		return
	}

	if typedError, ok := err.(entities.ErrServeRemovingEnv); ok {
		viewableError.Title = "Invalid sandbox state"
		viewableError.Message = fmt.Sprintf(
			"No port could be served in the sandbox \"%s\" because it's currently removing.\n\n"+
				"You must wait for the removing process to terminate.",
			config.ColorsBold(typedError.EnvName),
		)

		return
	}

	if typedError, ok := err.(entities.ErrUnserveCreatingEnv); ok {
		viewableError.Title = "Invalid sandbox state"
		viewableError.Message = fmt.Sprintf(
			"No port could be unserved in the sandbox \"%s\" because it's currently initializing.\n\n"+
				"You must wait for the initialization process to terminate.",
			config.ColorsBold(typedError.EnvName),
		)

		return
	}

	if typedError, ok := err.(entities.ErrUnserveRemovingEnv); ok {
		viewableError.Title = "Invalid sandbox state"
		viewableError.Message = fmt.Sprintf(
			"No port could be unserved in the sandbox \"%s\" because it's currently removing.\n\n"+
				"You must wait for the removing process to terminate.",
			config.ColorsBold(typedError.EnvName),
		)

		return
	}

	if _, ok := err.(entities.ErrUpdateInstanceTypeCreatingEnv); ok {
		viewableError.Title = "Unsupported operation"
		viewableError.Message = "The instance type cannot be updated once a sandbox is initializing."

		return
	}

	if typedError, ok := err.(entities.ErrInvalidPort); ok {
		viewableError.Title = "Invalid port"
		viewableError.Message = fmt.Sprintf(
			"The value \"%s\" is not a valid port.\n\n"+
				"Ports are numbers in the range 1-65535.",
			config.ColorsBold(typedError.InvalidPort),
		)

		return
	}

	if typedError, ok := err.(entities.ErrReservedPort); ok {
		viewableError.Title = "Reserved port"
		viewableError.Message = fmt.Sprintf(
			"The port \"%s\" is reserved for Eleven usage.",
			config.ColorsBold(typedError.ReservedPort),
		)

		return
	}

	if typedError, ok := err.(exceptions.ErrLoginError); ok {
		viewableError.Title = "GitHub authorization error"
		viewableError.Message = fmt.Sprintf(
			"An error has occured during the authorization of the Eleven application.\n\n%s",
			config.ColorsBold(typedError.Reason),
		)

		return
	}

	if typedError, ok := err.(exceptions.ErrMissingRequirements); ok {
		viewableError.Title = "Missing requirements"
		viewableError.Message = fmt.Sprintf(
			"The following requirements are missing:\n\n  - %s",
			config.ColorsBold(strings.Join(typedError.MissingRequirements, "\n\n  - ")),
		)

		return
	}

	if typedError, ok := err.(entities.ErrUnresolvableDomain); ok {
		viewableError.Title = "Unresolvable domain name"
		viewableError.Message = fmt.Sprintf(
			"The domain name \"%s\" doesn't resolve to \"%s\".\n\n"+
				"Please, add an A record and wait for DNS propagation.",
			config.ColorsBold(typedError.Domain),
			config.ColorsBold(typedError.EnvIPAddress),
		)

		return
	}

	if typedError, ok := err.(entities.ErrCloudflareSSLFull); ok {
		viewableError.Title = "Cloudflare SSL/TLS encryption mode set to \"Full\""
		viewableError.Message = fmt.Sprintf(
			"SSL/TLS encryption mode is set to \"Full\" for \"%s\".\n\n"+
				"Given that the SSL certificate for this domain is not yet issued, Clouflare refuses to connect to it.\n\n"+
				"Please, disable proxying for \"%s\", wait for DNS propagation and re-run the \"serve\" command.\n\n"+
				config.ColorsBold("Once served, proxying could be reactivated."),
			config.ColorsBold(typedError.Domain),
			config.ColorsBold(typedError.Domain),
		)

		return
	}

	if typedError, ok := err.(entities.ErrProxyForceHTTPS); ok {
		viewableError.Title = "HTTPS forced by proxy"
		viewableError.Message = fmt.Sprintf(
			"All requests to \"%s\" are redirected to \"%s\" by a proxy.\n\n"+
				"Given that the SSL certificate for this domain is not yet issued, all requests fail.\n\n"+
				"Please, disable proxying for \"%s\", wait for DNS propagation and re-run the \"serve\" command.\n\n"+
				config.ColorsBold("Once served, proxying could be reactivated."),
			config.ColorsBold("http://"+typedError.Domain),
			config.ColorsBold("https://"+typedError.Domain),
			config.ColorsBold(typedError.Domain),
		)

		return
	}

	if typedError, ok := err.(entities.ErrLetsEncryptTimedOut); ok {
		viewableError.Title = "Let's Encrypt timed out"
		viewableError.Message = fmt.Sprintf(
			"Let's Encrypt was unable to issue certificate for \"%s\" in the allocated time.\n\n"+
				"Please, wait a little bit and re-run the \"serve\" command.\n\n"+
				config.ColorsBold("Error: %s"),
			config.ColorsBold(typedError.Domain),
			config.ColorsBold(typedError.ReturnedError.Error()),
		)

		return
	}

	if typedError, ok := err.(entities.ErrInvalidDomain); ok {
		viewableError.Title = "Invalid domain name"
		viewableError.Message = fmt.Sprintf(
			"The value \"%s\" is not a valid domain name.",
			config.ColorsBold(typedError.Domain),
		)

		return
	}

	if typedError, ok := err.(entities.ErrEnvCloudInitError); ok {
		viewableError.Title = "Cloud init error"
		viewableError.Logs = typedError.Logs
		viewableError.Message = typedError.ErrorMessage

		return
	}

	if typedError, ok := err.(entities.ErrInvalidEnvName); ok {
		viewableError.Title = "Invalid sandbox name"
		viewableError.Message = fmt.Sprintf(
			"The string \"%s\" is not a valid sandbox name.\n\n"+
				"Sandbox names must match \"%s\" and be less than or equal to %s characters.",
			config.ColorsBold(typedError.EnvName),
			config.ColorsBold(typedError.EnvNameRegExp),
			config.ColorsBold(fmt.Sprintf("%d", typedError.EnvNameMaxLength)),
		)

		return
	}

	bold := config.ColorsBold

	if githubErr, ok := err.(*github.ErrorResponse); ok {
		viewableError.Title = "GitHub API error"
		viewableError.Message = fmt.Sprintf(
			"%s\n\nAn error occurred while calling the GitHub API.\n\n"+
				"You could try to fix it (using the details above) or open a new issue at: https://github.com/eleven-sh/cli/issues/new",
			bold(githubErr.Error()),
		)

		return
	}

	if status, ok := status.FromError(err); ok {
		viewableError.Title = "Eleven agent error"

		errorMessage := status.Message()

		if len(errorMessage) >= 2 {
			errorMessage = strings.ToTitle(errorMessage[0:1]) + errorMessage[1:] + "."
		}

		viewableError.Message = fmt.Sprintf(
			"%s\n\nAn error occured in the Eleven agent.\n\n"+
				"You could try to fix it (using the details above) or open a new issue at: https://github.com/eleven-sh/cli/issues/new",
			bold(errorMessage),
		)

		if status.Code() != codes.Unknown {
			viewableError.Message += "\n\n" +
				bold("Error code: ") +
				status.Code().String()
		}

		return
	}

	viewableError.Title = "Unknown error"
	viewableError.Message = fmt.Sprintf(
		"%s\n\nAn unknown error occurred.\n\n"+
			"You could try to fix it (using the details above) or open a new issue at: https://github.com/eleven-sh/cli/issues/new",
		bold(err.Error()),
	)

	return
}
