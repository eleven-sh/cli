package presenters

import (
	"github.com/eleven-sh/cli/internal/features"
	"github.com/eleven-sh/cli/internal/views"
)

type RemoveViewer interface {
	View(views.RemoveViewData)
}

type RemovePresenter struct {
	viewableErrorBuilder ViewableErrorBuilder
	viewer               RemoveViewer
}

func NewRemovePresenter(
	viewableErrorBuilder ViewableErrorBuilder,
	viewer RemoveViewer,
) RemovePresenter {

	return RemovePresenter{
		viewableErrorBuilder: viewableErrorBuilder,
		viewer:               viewer,
	}
}

func (r RemovePresenter) PresentToView(response features.RemoveResponse) {
	viewData := views.RemoveViewData{}

	if response.Error == nil {
		envName := response.Content.EnvName

		viewData.Content = views.RemoveViewDataContent{
			Message: "The sandbox \"" + envName + "\" was removed.",
		}

		r.viewer.View(viewData)
		return
	}

	viewData.Error = r.viewableErrorBuilder.Build(response.Error)
	r.viewer.View(viewData)
}
