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

func ProvideHetznerUninstallFeature(elevenConfigDir, region, context string) features.UninstallFeature {
	return provideHetznerUninstallFeature(
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

func provideHetznerUninstallFeature(
	userConfigEnvVarsResolverOpts hetznerProviderUserConfig.EnvVarsResolverOpts,
	userConfigFilesResolverOpts hetznerProviderUserConfig.FilesResolverOpts,
	userConfigLocalResolverOpts hetznerCLI.UserConfigLocalResolverOpts,
	serviceBuilderOpts hetznerProviderService.BuilderOpts,
) features.UninstallFeature {
	panic(
		wire.Build(
			viewSet,
			hetznerServiceBuilderSet,
			hetznerViewableErrorBuilder,

			stepperSet,

			wire.Bind(new(features.UninstallOutputHandler), new(featuresCLI.UninstallOutputHandler)),
			featuresCLI.NewUninstallOutputHandler,

			wire.Bind(new(featuresCLI.UninstallPresenter), new(presenters.UninstallPresenter)),
			presenters.NewUninstallPresenter,

			wire.Bind(new(presenters.UninstallViewer), new(views.UninstallView)),
			views.NewUninstallView,

			features.NewUninstallFeature,
		),
	)
}
