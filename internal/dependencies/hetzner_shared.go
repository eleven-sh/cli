// go:build wireinject
//go:build wireinject
// +build wireinject

package dependencies

import (
	hetznerCLI "github.com/eleven-sh/cli/internal/cloudproviders/hetzner"
	"github.com/eleven-sh/cli/internal/presenters"
	"github.com/eleven-sh/cli/internal/system"
	"github.com/eleven-sh/eleven/entities"
	hetznerProviderConfig "github.com/eleven-sh/hetzner-cloud-provider/config"
	hetznerProviderService "github.com/eleven-sh/hetzner-cloud-provider/service"
	hetznerProviderUserConfig "github.com/eleven-sh/hetzner-cloud-provider/userconfig"
	"github.com/google/wire"
)

var hetznerViewableErrorBuilder = wire.NewSet(
	wire.Bind(new(presenters.ViewableErrorBuilder), new(hetznerCLI.HetznerViewableErrorBuilder)),
	hetznerCLI.NewHetznerViewableErrorBuilder,
)

var hetznerServiceBuilderSet = wire.NewSet(
	wire.Bind(new(hetznerProviderUserConfig.ContextLoader), new(hetznerProviderConfig.ContextLoader)),
	hetznerProviderConfig.NewContextLoader,

	wire.Bind(new(hetznerCLI.UserConfigFilesResolver), new(hetznerProviderUserConfig.FilesResolver)),
	hetznerProviderUserConfig.NewFilesResolver,

	wire.Bind(new(hetznerProviderUserConfig.EnvVarsGetter), new(system.EnvVars)),
	system.NewEnvVars,

	wire.Bind(new(hetznerCLI.UserConfigEnvVarsResolver), new(hetznerProviderUserConfig.EnvVarsResolver)),
	hetznerProviderUserConfig.NewEnvVarsResolver,

	wire.Bind(new(hetznerProviderService.UserConfigResolver), new(hetznerCLI.UserConfigLocalResolver)),
	hetznerCLI.NewUserConfigLocalResolver,

	wire.Bind(new(hetznerProviderService.UserConfigValidator), new(hetznerProviderConfig.UserConfigValidator)),
	hetznerProviderConfig.NewUserConfigValidator,

	wire.Bind(new(entities.CloudServiceBuilder), new(hetznerProviderService.Builder)),
	hetznerProviderService.NewBuilder,
)
