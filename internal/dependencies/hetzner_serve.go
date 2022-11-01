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

func ProvideHetznerServeFeature(elevenConfigDir, region, context string) features.ServeFeature {
	return provideHetznerServeFeature(
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

func provideHetznerServeFeature(
	userConfigEnvVarsResolverOpts hetznerProviderUserConfig.EnvVarsResolverOpts,
	userConfigFilesResolverOpts hetznerProviderUserConfig.FilesResolverOpts,
	userConfigLocalResolverOpts hetznerCLI.UserConfigLocalResolverOpts,
	serviceBuilderOpts hetznerProviderService.BuilderOpts,
) features.ServeFeature {
	panic(
		wire.Build(
			viewSet,
			hetznerServiceBuilderSet,
			hetznerViewableErrorBuilder,

			loggerSet,

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
