// go:build wireinject
//go:build wireinject
// +build wireinject

package dependencies

import (
	"github.com/eleven-sh/cli/internal/features"
	"github.com/eleven-sh/cli/internal/presenters"
	"github.com/eleven-sh/cli/internal/views"
	"github.com/google/wire"
)

func ProvideLoginFeature() features.LoginFeature {
	panic(
		wire.Build(
			viewSet,
			elevenViewableErrorBuilder,

			loggerSet,

			browserManagerSet,

			userConfigManagerSet,

			sleeperSet,

			githubManagerSet,

			wire.Bind(new(features.LoginPresenter), new(presenters.LoginPresenter)),
			presenters.NewLoginPresenter,

			wire.Bind(new(presenters.LoginViewer), new(views.LoginView)),
			views.NewLoginView,

			features.NewLoginFeature,
		),
	)
}
