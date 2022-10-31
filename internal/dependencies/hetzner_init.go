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

func ProvideHetznerInitFeature(elevenConfigDir, region, context string) features.InitFeature {
	return provideHetznerInitFeature(
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

func provideHetznerInitFeature(
	userConfigEnvVarsResolverOpts hetznerProviderUserConfig.EnvVarsResolverOpts,
	userConfigFilesResolverOpts hetznerProviderUserConfig.FilesResolverOpts,
	userConfigLocalResolverOpts hetznerCLI.UserConfigLocalResolverOpts,
	serviceBuilderOpts hetznerProviderService.BuilderOpts,
) features.InitFeature {
	panic(
		wire.Build(
			viewSet,
			hetznerServiceBuilderSet,
			hetznerViewableErrorBuilder,

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
