package features

import (
	"encoding/json"
	"io"
	"strconv"

	agentConfig "github.com/eleven-sh/agent/config"
	"github.com/eleven-sh/agent/proto"
	"github.com/eleven-sh/cli/internal/agent"
	"github.com/eleven-sh/cli/internal/config"

	cliEntities "github.com/eleven-sh/cli/internal/entities"
	"github.com/eleven-sh/cli/internal/interfaces"

	"github.com/eleven-sh/eleven/actions"
	"github.com/eleven-sh/eleven/entities"
	"github.com/eleven-sh/eleven/features"
)

type InitResponse struct {
	Error   error
	Content InitResponseContent
}

type InitResponseContent struct {
	EnvName                   string
	EnvLocalSSHConfigHostname string
	EnvPublicIPAddress        string
	EnvAlreadyCreated         bool
	EnvRuntimes               entities.EnvRuntimes
}

type InitPresenter interface {
	PresentToView(InitResponse)
}

type InitOutputHandler struct {
	userConfig         interfaces.UserConfigManager
	presenter          InitPresenter
	agentClientBuilder agent.ClientBuilder
	github             interfaces.GitHubManager
	logger             interfaces.Logger
	sshConfig          interfaces.SSHConfigManager
	sshKeys            interfaces.SSHKeysManager
	sshKnownHosts      interfaces.SSHKnownHostsManager
}

func NewInitOutputHandler(
	userConfig interfaces.UserConfigManager,
	presenter InitPresenter,
	agentClientBuilder agent.ClientBuilder,
	github interfaces.GitHubManager,
	logger interfaces.Logger,
	sshConfig interfaces.SSHConfigManager,
	sshKeys interfaces.SSHKeysManager,
	sshKnownHosts interfaces.SSHKnownHostsManager,
) InitOutputHandler {

	return InitOutputHandler{
		userConfig:         userConfig,
		presenter:          presenter,
		agentClientBuilder: agentClientBuilder,
		github:             github,
		logger:             logger,
		sshConfig:          sshConfig,
		sshKnownHosts:      sshKnownHosts,
		sshKeys:            sshKeys,
	}
}

func (i InitOutputHandler) HandleOutput(output features.InitOutput) error {
	stepper := output.Stepper

	handleError := func(err error) error {
		stepper.StopCurrentStep()

		i.presenter.PresentToView(InitResponse{
			Error: err,
		})

		return err
	}

	if output.Error != nil {
		return handleError(output.Error)
	}

	env := output.Content.Env
	envCreated := output.Content.EnvCreated
	envAlreadyCreated := !envCreated

	var envAdditionalProperties *cliEntities.EnvAdditionalProperties

	if len(env.AdditionalPropertiesJSON) > 0 {
		err := json.Unmarshal(
			[]byte(env.AdditionalPropertiesJSON),
			&envAdditionalProperties,
		)

		if err != nil {
			return handleError(err)
		}
	}

	if envAdditionalProperties == nil {
		envAdditionalProperties = &cliEntities.EnvAdditionalProperties{}
	}

	if envCreated {
		stepper.StartTemporaryStep(
			"Building the sandbox",
		)

		agentClient := i.agentClientBuilder.Build(
			agent.NewDefaultClientConfig(
				[]byte(env.SSHKeyPairPEMContent),
				env.InstancePublicIPAddress,
			),
		)

		err := agentClient.InitInstance(
			&proto.InitInstanceRequest{
				EnvName:         env.Name,
				EnvNameSlug:     env.GetNameSlug(),
				EnvRepos:        cliEntities.BuildProtoRepositoriesFromEnv(env),
				GithubUserEmail: i.userConfig.GetString(config.UserConfigKeyGitHubEmail),
				UserFullName:    i.userConfig.GetString(config.UserConfigKeyGitHubFullName),
			},
			func(stream agent.InitInstanceStream) error {

				logs := ""

				for {
					reply, err := stream.Recv()

					if err == io.EOF {
						break
					}

					if err != nil {
						stepper.StopCurrentStep()
						i.logger.Log(logs)
						return err
					}

					if reply.GithubSshPublicKeyContent != nil &&
						envAdditionalProperties.GitHubCreatedSSHKeyId == nil {

						sshKeyInGitHub, err := i.github.CreateSSHKey(
							i.userConfig.GetString(config.UserConfigKeyGitHubAccessToken),
							env.GetSSHKeyPairName(),
							reply.GetGithubSshPublicKeyContent(),
						)

						if err != nil {
							return err
						}

						envAdditionalProperties.GitHubCreatedSSHKeyId = sshKeyInGitHub.ID
						err = env.SetAdditionalPropertiesJSON(envAdditionalProperties)

						if err != nil {
							return err
						}

						err = actions.UpdateEnvInConfig(
							stepper,
							output.Content.CloudService,
							output.Content.ElevenConfig,
							output.Content.Cluster,
							env,
						)

						if err != nil {
							return err
						}
					}

					if len(reply.LogLine) > 0 {
						logs += reply.LogLine
					}
				}

				return nil
			},
		)

		if err != nil {
			return handleError(err)
		}

		// Remove default Caddy config (Welcome page)
		err = cliEntities.ReconcileServedPortsState(
			env,
			i.agentClientBuilder,
		)

		if err != nil {
			return handleError(err)
		}

		err = agentClient.InstallRuntimes(
			&proto.InstallRuntimesRequest{
				Runtimes: output.Content.Runtimes,
			},
			func(stream agent.InstallRuntimesStream) error {

				logs := ""

				for {
					reply, err := stream.Recv()

					if err == io.EOF {
						break
					}

					if err != nil {
						stepper.StopCurrentStep()
						i.logger.Log(logs)
						return err
					}

					if len(reply.LogLineHeader) > 0 {
						stepper.StartTemporaryStep(
							reply.LogLineHeader,
						)
						logs = ""
					}

					if len(reply.LogLine) > 0 {
						logs += reply.LogLine
					}
				}

				return nil
			},
		)

		if err != nil {
			return handleError(err)
		}
	}

	stepper.StartTemporaryStepWithoutNewLine(
		"Updating your local SSH configuration",
	)

	sshPEMPath, err := i.sshKeys.CreateOrReplacePEM(
		env.GetSSHKeyPairName(),
		env.SSHKeyPairPEMContent,
	)

	if err != nil {
		return handleError(err)
	}

	sshServerListenPort, err := strconv.ParseInt(
		agentConfig.SSHServerListenPort,
		10,
		64,
	)

	if err != nil {
		return handleError(err)
	}

	err = i.sshConfig.AddOrReplaceHost(
		env.LocalSSHConfigHostname,
		env.InstancePublicIPAddress,
		sshPEMPath,
		agentConfig.ElevenUserName,
		sshServerListenPort,
	)

	if err != nil {
		return handleError(err)
	}

	for _, sshHostKey := range env.SSHHostKeys {
		err := i.sshKnownHosts.AddOrReplace(
			env.InstancePublicIPAddress,
			sshHostKey.Algorithm,
			sshHostKey.Fingerprint,
		)

		if err != nil {
			return handleError(err)
		}
	}

	if !envAlreadyCreated {
		err := output.Content.SetEnvAsCreated()

		if err != nil {
			return handleError(err)
		}
	}

	stepper.StopCurrentStep()

	i.presenter.PresentToView(InitResponse{
		Content: InitResponseContent{
			EnvName:                   env.Name,
			EnvLocalSSHConfigHostname: env.LocalSSHConfigHostname,
			EnvPublicIPAddress:        env.InstancePublicIPAddress,
			EnvAlreadyCreated:         envAlreadyCreated,
			EnvRuntimes:               env.Runtimes,
		},
	})

	return nil
}
