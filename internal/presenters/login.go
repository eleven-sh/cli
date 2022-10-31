package presenters

import (
	"github.com/eleven-sh/cli/internal/features"
	"github.com/eleven-sh/cli/internal/views"
)

type LoginViewer interface {
	View(views.LoginViewData)
}

type LoginPresenter struct {
	viewableErrorBuilder ViewableErrorBuilder
	viewer               LoginViewer
}

func NewLoginPresenter(
	viewableErrorBuilder ViewableErrorBuilder,
	viewer LoginViewer,
) LoginPresenter {

	return LoginPresenter{
		viewableErrorBuilder: viewableErrorBuilder,
		viewer:               viewer,
	}
}

func (l LoginPresenter) PresentToView(response features.LoginResponse) {
	viewData := views.LoginViewData{}

	if response.Error == nil {
		viewData.Content = views.LoginViewDataContent{
			Message: "Your GitHub account is now connected.",
		}

		l.viewer.View(viewData)
		return
	}

	viewData.Error = l.viewableErrorBuilder.Build(response.Error)
	l.viewer.View(viewData)
}
