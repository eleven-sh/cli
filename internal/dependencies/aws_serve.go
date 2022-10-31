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

func ProvideAWSServeFeature(region, profile, credentialsFilePath, configFilePath string) features.ServeFeature {
	return provideAWSServeFeature(
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

func provideAWSServeFeature(
	userConfigEnvVarsResolverOpts awsProviderUserConfig.EnvVarsResolverOpts,
	userConfigFilesResolverOpts awsProviderUserConfig.FilesResolverOpts,
	userConfigLocalResolverOpts awsCLI.UserConfigLocalResolverOpts,
) features.ServeFeature {
	panic(
		wire.Build(
			viewSet,
			awsServiceBuilderSet,
			awsViewableErrorBuilder,

			stepperSet,

			agentSet,

			wire.Bind(new(features.ServeOutputHandler), new(featuresCLI.ServeOutputHandler)),
			featuresCLI.NewServeOutputHandler,

			wire.Bind(new(featuresCLI.ServePresenter), new(presenters.ServePresenter)),
			presenters.NewServePresenter,

			wire.Bind(new(presenters.ServeViewer), new(views.ServeView)),
			views.NewServeView,

			features.NewServeFeature,
		),
	)
}
