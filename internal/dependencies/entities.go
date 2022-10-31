// go:build wireinject
//go:build wireinject
// +build wireinject

package dependencies

import (
	"github.com/eleven-sh/cli/internal/entities"
	"github.com/google/wire"
)

func ProvideEnvRepositoriesResolver() entities.EnvRepositoriesResolver {
	panic(
		wire.Build(
			loggerSet,

			userConfigManagerSet,

			githubManagerSet,

			entities.NewEnvRepositoriesResolver,
		),
	)
}
