// go:build wireinject
//go:build wireinject
// +build wireinject

package dependencies

import (
	hetznerCLI "github.com/eleven-sh/cli/internal/cloudproviders/hetzner"
	featuresCLI "github.com/eleven-sh/cli/internal/features"
	"github.com/eleven-sh/cli/internal/presenters"
	"github.com/eleven-sh/cli/internal/views"
	"github.com/eleven-sh/eleven/features"
	hetznerProviderService "github.com/eleven-sh/hetzner-cloud-provider/service"
	hetznerProviderUserConfig "github.com/eleven-sh/hetzner-cloud-provider/userconfig"
	"github.com/google/wire"
)

func ProvideHetznerEditFeature(elevenConfigDir, region, context string) features.EditFeature {
	return provideHetznerEditFeature(
		hetznerProviderUserConfig.EnvVarsResolverOpts{
			Region: region,
		},

		hetznerProviderUserConfig.FilesResolverOpts{
			Region:  region,
			Context: context,
		},

		hetznerCLI.UserConfigLocalResolverOpts{
			Context: context,
		},

		hetznerProviderService.BuilderOpts{
			ElevenConfigDir: elevenConfigDir,
		},
	)
}

func provideHetznerEditFeature(
	userConfigEnvVarsResolverOpts hetznerProviderUserConfig.EnvVarsResolverOpts,
	userConfigFilesResolverOpts hetznerProviderUserConfig.FilesResolverOpts,
	userConfigLocalResolverOpts hetznerCLI.UserConfigLocalResolverOpts,
	serviceBuilderOpts hetznerProviderService.BuilderOpts,
) features.EditFeature {
	panic(
		wire.Build(
			viewSet,
			hetznerServiceBuilderSet,
			hetznerViewableErrorBuilder,

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
