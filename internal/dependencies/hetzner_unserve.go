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

func ProvideHetznerUnserveFeature(elevenConfigDir, region, context string) features.UnserveFeature {
	return provideHetznerUnserveFeature(
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

func provideHetznerUnserveFeature(
	userConfigEnvVarsResolverOpts hetznerProviderUserConfig.EnvVarsResolverOpts,
	userConfigFilesResolverOpts hetznerProviderUserConfig.FilesResolverOpts,
	userConfigLocalResolverOpts hetznerCLI.UserConfigLocalResolverOpts,
	serviceBuilderOpts hetznerProviderService.BuilderOpts,
) features.UnserveFeature {
	panic(
		wire.Build(
			viewSet,
			hetznerServiceBuilderSet,
			hetznerViewableErrorBuilder,

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
