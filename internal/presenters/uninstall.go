package presenters

import (
	"fmt"

	"github.com/eleven-sh/cli/internal/config"
	"github.com/eleven-sh/cli/internal/features"
	"github.com/eleven-sh/cli/internal/views"
)

type UninstallViewer interface {
	View(views.UninstallViewData)
}

type UninstallPresenter struct {
	viewableErrorBuilder ViewableErrorBuilder
	viewer               UninstallViewer
}

func NewUninstallPresenter(
	viewableErrorBuilder ViewableErrorBuilder,
	viewer UninstallViewer,
) UninstallPresenter {

	return UninstallPresenter{
		viewableErrorBuilder: viewableErrorBuilder,
		viewer:               viewer,
	}
}

func (u UninstallPresenter) PresentToView(response features.UninstallResponse) {
	viewData := views.UninstallViewData{}

	if response.Error == nil {
		bold := config.ColorsBold

		elevenAlreadyUninstalled := response.Content.ElevenAlreadyUninstalled

		viewDataMessage := response.Content.SuccessMessage
		if elevenAlreadyUninstalled {
			viewDataMessage = response.Content.AlreadyUninstalledMessage
		}

		viewDataSubtext := fmt.Sprintf(
			"If you want to remove Eleven entirely:\n\n"+
				"  - Remove the Eleven CLI (located at %s)\n\n"+
				"  - Remove the Eleven configuration (located at %s)\n\n"+
				"  - Unauthorize the Eleven application on GitHub by going to: %s",
			bold(response.Content.ElevenExecutablePath),
			bold(response.Content.ElevenConfigDirPath),
			bold("https://github.com/settings/applications"),
		)

		viewData.Content = views.UninstallViewDataContent{
			ShowAsWarning: elevenAlreadyUninstalled,
			Message:       viewDataMessage,
			Subtext:       viewDataSubtext,
		}

		u.viewer.View(viewData)
		return
	}

	viewData.Error = u.viewableErrorBuilder.Build(response.Error)
	u.viewer.View(viewData)
}
