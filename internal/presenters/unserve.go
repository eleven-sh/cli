package presenters

import (
	"fmt"

	"github.com/eleven-sh/cli/internal/features"
	"github.com/eleven-sh/cli/internal/views"
)

type UnserveViewer interface {
	View(views.UnserveViewData)
}

type UnservePresenter struct {
	viewableErrorBuilder ViewableErrorBuilder
	viewer               UnserveViewer
}

func NewUnservePresenter(
	viewableErrorBuilder ViewableErrorBuilder,
	viewer UnserveViewer,
) UnservePresenter {

	return UnservePresenter{
		viewableErrorBuilder: viewableErrorBuilder,
		viewer:               viewer,
	}
}

func (u UnservePresenter) PresentToView(response features.UnserveResponse) {
	viewData := views.UnserveViewData{}

	if response.Error == nil {
		unservedPort := response.Content.Port

		viewDataMessage := fmt.Sprintf(
			"The port \"%s\" is now unreachable from outside.",
			unservedPort,
		)

		viewData.Content = views.UnserveViewDataContent{
			Message: viewDataMessage,
		}

		u.viewer.View(viewData)
		return
	}

	viewData.Error = u.viewableErrorBuilder.Build(response.Error)
	u.viewer.View(viewData)
}
