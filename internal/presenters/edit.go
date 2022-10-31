package presenters

import (
	"github.com/eleven-sh/cli/internal/features"
	"github.com/eleven-sh/cli/internal/views"
)

type EditViewer interface {
	View(views.EditViewData)
}

type EditPresenter struct {
	viewableErrorBuilder ViewableErrorBuilder
	viewer               EditViewer
}

func NewEditPresenter(
	viewableErrorBuilder ViewableErrorBuilder,
	viewer EditViewer,
) EditPresenter {

	return EditPresenter{
		viewableErrorBuilder: viewableErrorBuilder,
		viewer:               viewer,
	}
}

func (e EditPresenter) PresentToView(response features.EditResponse) {
	viewData := views.EditViewData{}

	if response.Error == nil {
		viewData.Content = views.EditViewDataContent{
			Message: "Your editor is now open and connected to your sandbox.",
		}

		e.viewer.View(viewData)
		return
	}

	viewData.Error = e.viewableErrorBuilder.Build(response.Error)
	e.viewer.View(viewData)
}
