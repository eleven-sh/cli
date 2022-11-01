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

func ProvideAWSEditFeature(region, profile, credentialsFilePath, configFilePath string) features.EditFeature {
	return provideAWSEditFeature(
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

func provideAWSEditFeature(
	userConfigEnvVarsResolverOpts awsProviderUserConfig.EnvVarsResolverOpts,
	userConfigFilesResolverOpts awsProviderUserConfig.FilesResolverOpts,
	userConfigLocalResolverOpts awsCLI.UserConfigLocalResolverOpts,
) features.EditFeature {
	panic(
		wire.Build(
			viewSet,
			awsServiceBuilderSet,
			awsViewableErrorBuilder,

			loggerSet,

			stepperSet,

			vscodeProcessManagerSet,

			vscodeExtensionsManagerSet,

			wire.Bind(new(features.EditOutputHandler), new(featuresCLI.EditOutputHandler)),
			featuresCLI.NewEditOutputHandler,

			wire.Bind(new(featuresCLI.EditPresenter), new(presenters.EditPresenter)),
			presenters.NewEditPresenter,

			wire.Bind(new(presenters.EditViewer), new(views.EditView)),
			views.NewEditView,

			features.NewEditFeature,
		),
	)
}
