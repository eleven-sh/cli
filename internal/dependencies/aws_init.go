// go:build wireinject
//go:build wireinject
// +build wireinject

package dependencies

import (
	awsProviderUserConfig "github.com/eleven-sh/aws-cloud-provider/userconfig"
	awsCLI "github.com/eleven-sh/cli/internal/cloudproviders/aws"
	featuresCLI "github.com/eleven-sh/cli/internal/features"
	"github.com/eleven-sh/cli/internal/presenters"
	"github.com/eleven-sh/cli/internal/views"
	"github.com/eleven-sh/eleven/features"
	"github.com/google/wire"
)

func ProvideAWSInitFeature(region, profile, credentialsFilePath, configFilePath string) features.InitFeature {
	return provideAWSInitFeature(
		awsProviderUserConfig.EnvVarsResolverOpts{
			Region: region,
		},

		awsProviderUserConfig.FilesResolverOpts{
			Region:              region,
			Profile:             profile,
			CredentialsFilePath: credentialsFilePath,
			ConfigFilePath:      configFilePath,
		},

		awsCLI.UserConfigLocalResolverOpts{
			Profile: profile,
		},
	)
}

func provideAWSInitFeature(
	userConfigEnvVarsResolverOpts awsProviderUserConfig.EnvVarsResolverOpts,
	userConfigFilesResolverOpts awsProviderUserConfig.FilesResolverOpts,
	userConfigLocalResolverOpts awsCLI.UserConfigLocalResolverOpts,
) features.InitFeature {
	panic(
		wire.Build(
			viewSet,
			awsServiceBuilderSet,
			awsViewableErrorBuilder,

			userConfigManagerSet,

			agentSet,

			githubManagerSet,

			loggerSet,

			sshConfigManagerSet,

			sshKnownHostsManagerSet,

			sshKeysManagerSet,

			stepperSet,

			wire.Bind(new(features.InitOutputHandler), new(featuresCLI.InitOutputHandler)),
			featuresCLI.NewInitOutputHandler,

			wire.Bind(new(featuresCLI.InitPresenter), new(presenters.InitPresenter)),
			presenters.NewInitPresenter,

			wire.Bind(new(presenters.InitViewer), new(views.InitView)),
			views.NewInitView,

			features.NewInitFeature,
		),
	)
}
