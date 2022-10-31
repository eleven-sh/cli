package hooks

import (
	"encoding/json"

	"github.com/eleven-sh/cli/internal/config"
	cliEntities "github.com/eleven-sh/cli/internal/entities"
	"github.com/eleven-sh/cli/internal/interfaces"
	"github.com/eleven-sh/eleven/entities"
)

type PreRemove struct {
	sshConfig     interfaces.SSHConfigManager
	sshKeys       interfaces.SSHKeysManager
	sshKnownHosts interfaces.SSHKnownHostsManager
	userConfig    interfaces.UserConfigManager
	github        interfaces.GitHubManager
}

func NewPreRemove(
	sshConfig interfaces.SSHConfigManager,
	sshKeys interfaces.SSHKeysManager,
	sshKnownHosts interfaces.SSHKnownHostsManager,
	userConfig interfaces.UserConfigManager,
	github interfaces.GitHubManager,
) PreRemove {

	return PreRemove{
		sshConfig:     sshConfig,
		sshKeys:       sshKeys,
		sshKnownHosts: sshKnownHosts,
		userConfig:    userConfig,
		github:        github,
	}
}

func (p PreRemove) Run(
	cloudService entities.CloudService,
	elevenConfig *entities.Config,
	cluster *entities.Cluster,
	env *entities.Env,
) error {

	err := p.sshKeys.RemovePEMIfExists(env.GetSSHKeyPairName())

	if err != nil {
		return err
	}

	err = p.sshConfig.RemoveHostIfExists(env.LocalSSHConfigHostname)

	if err != nil {
		return err
	}

	sshHostname := env.InstancePublicIPAddress
	err = p.sshKnownHosts.RemoveIfExists(sshHostname)

	if err != nil {
		return err
	}

	// User could remove dev env in creating state
	// (in case of error for example)
	if len(env.AdditionalPropertiesJSON) == 0 {
		return nil
	}

	var envAdditionalProperties *cliEntities.EnvAdditionalProperties
	err = json.Unmarshal(
		[]byte(env.AdditionalPropertiesJSON),
		&envAdditionalProperties,
	)

	if err != nil {
		return err
	}

	githubAccessToken := p.userConfig.GetString(
		config.UserConfigKeyGitHubAccessToken,
	)

	if envAdditionalProperties.GitHubCreatedSSHKeyId != nil {
		err = p.github.RemoveSSHKey(
			githubAccessToken,
			*envAdditionalProperties.GitHubCreatedSSHKeyId,
		)

		if err != nil && !p.github.IsNotFoundError(err) {
			return err
		}
	}

	return nil
}
