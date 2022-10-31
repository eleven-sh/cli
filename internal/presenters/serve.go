package presenters

import (
	"fmt"

	"github.com/asaskevich/govalidator"
	"github.com/eleven-sh/cli/internal/config"
	"github.com/eleven-sh/cli/internal/features"
	"github.com/eleven-sh/cli/internal/views"
)

type ServeViewer interface {
	View(views.ServeViewData)
}

type ServePresenter struct {
	viewableErrorBuilder ViewableErrorBuilder
	viewer               ServeViewer
}

func NewServePresenter(
	viewableErrorBuilder ViewableErrorBuilder,
	viewer ServeViewer,
) ServePresenter {

	return ServePresenter{
		viewableErrorBuilder: viewableErrorBuilder,
		viewer:               viewer,
	}
}

func (s ServePresenter) PresentToView(response features.ServeResponse) {
	viewData := views.ServeViewData{}

	if response.Error == nil {
		portBinding := response.Content.PortBinding

		envIPAddress := response.Content.EnvPublicIPAddress
		servedPort := response.Content.Port

		portReachableAt := envIPAddress + ":" + portBinding
		if !govalidator.IsPort(portBinding) {
			portReachableAt = "https://" + portBinding
		}

		viewDataMessage := fmt.Sprintf(
			"The port \"%s\" is now reachable at: %s",
			servedPort,
			config.ColorsBlue(portReachableAt),
		)

		viewData.Content = views.ServeViewDataContent{
			Message: viewDataMessage,
		}

		s.viewer.View(viewData)
		return
	}

	viewData.Error = s.viewableErrorBuilder.Build(response.Error)
	s.viewer.View(viewData)
}
