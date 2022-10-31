package hetzner

import (
	"errors"

	"github.com/eleven-sh/hetzner-cloud-provider/userconfig"
)

type UserConfigLocalResolverOpts struct {
	Context string
}

//go:generate go run github.com/golang/mock/mockgen -destination=../../mocks/hetzner_user_config_env_vars_resolver.go -package=mocks -mock_names UserConfigEnvVarsResolver=HetznerUserConfigEnvVarsResolver github.com/eleven-sh/cli/internal/cloudproviders/hetzner UserConfigEnvVarsResolver
type UserConfigEnvVarsResolver interface {
	Resolve() (*userconfig.Config, error)
}

//go:generate go run github.com/golang/mock/mockgen -destination=../../mocks/hetzner_user_config_files_resolver.go -package=mocks -mock_names UserConfigFilesResolver=HetznerUserConfigFilesResolver github.com/eleven-sh/cli/internal/cloudproviders/hetzner UserConfigFilesResolver
type UserConfigFilesResolver interface {
	Resolve() (*userconfig.Config, error)
}

// UserConfigLocalResolver represents the default implementation
// of the UserConfigResolver interface, used by most Hetzner commands via
// the SDKConfigStaticBuilder.
//
// It retrieves the Hetzner account configuration from environment variables
// (via the UserConfigLocalEnvVarsResolver interface) and fallback to config
// files (via the UserConfigLocalFilesResolver interface) otherwise.
//
type UserConfigLocalResolver struct {
	envVarsResolver     UserConfigEnvVarsResolver
	configFilesResolver UserConfigFilesResolver
	opts                UserConfigLocalResolverOpts
}

// NewUserConfigLocalResolver constructs
// the UserConfigLocalResolver struct.
// Used by Wire in dependencies.
//
func NewUserConfigLocalResolver(
	envVarsResolver UserConfigEnvVarsResolver,
	configFilesResolver UserConfigFilesResolver,
	opts UserConfigLocalResolverOpts,
) UserConfigLocalResolver {

	return UserConfigLocalResolver{
		envVarsResolver:     envVarsResolver,
		configFilesResolver: configFilesResolver,
		opts:                opts,
	}
}

// Resolve retrieves the Hetzner account configuration from environment variables
// and fallback to config files if no environment variables were found.
//
// If the Profile option is set, environment variables are ignored
// and the profile is directly loaded from config files.
//
func (u UserConfigLocalResolver) Resolve() (*userconfig.Config, error) {
	var userConfig *userconfig.Config
	var err error

	if len(u.opts.Context) == 0 {
		userConfig, err = u.envVarsResolver.Resolve()

		if err != nil && !errors.Is(err, userconfig.ErrMissingConfig) {
			return nil, err
		}
	}

	if userConfig == nil {
		userConfig, err = u.configFilesResolver.Resolve()

		if err != nil {
			return nil, err
		}
	}

	return userConfig, nil
}
