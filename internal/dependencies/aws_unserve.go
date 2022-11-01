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

func ProvideAWSUnserveFeature(region, profile, credentialsFilePath, configFilePath string) features.UnserveFeature {
	return provideAWSUnserveFeature(
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

func provideAWSUnserveFeature(
	userConfigEnvVarsResolverOpts awsProviderUserConfig.EnvVarsResolverOpts,
	userConfigFilesResolverOpts awsProviderUserConfig.FilesResolverOpts,
	userConfigLocalResolverOpts awsCLI.UserConfigLocalResolverOpts,
) features.UnserveFeature {
	panic(
		wire.Build(
			viewSet,
			awsServiceBuilderSet,
			awsViewableErrorBuilder,

			loggerSet,

			stepperSet,

			agentSet,

			wire.Bind(new(features.UnserveOutputHandler), new(featuresCLI.UnserveOutputHandler)),
			featuresCLI.NewUnserveOutputHandler,

			wire.Bind(new(featuresCLI.UnservePresenter), new(presenters.UnservePresenter)),
			presenters.NewUnservePresenter,

			wire.Bind(new(presenters.UnserveViewer), new(views.UnserveView)),
			views.NewUnserveView,

			features.NewUnserveFeature,
		),
	)
}
