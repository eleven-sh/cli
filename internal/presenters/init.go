package presenters

import (
	"fmt"

	"github.com/eleven-sh/cli/internal/config"
	"github.com/eleven-sh/cli/internal/features"
	"github.com/eleven-sh/cli/internal/globals"
	"github.com/eleven-sh/cli/internal/views"
)

//go:generate go run github.com/golang/mock/mockgen -destination=../mocks/presenters_init.go -package=mocks -mock_names InitViewer=PresentersInitViewer github.com/eleven-sh/cli/internal/presenters InitViewer
type InitViewer interface {
	View(views.InitViewData)
}

type InitPresenter struct {
	viewableErrorBuilder ViewableErrorBuilder
	viewer               InitViewer
}

func NewInitPresenter(
	viewableErrorBuilder ViewableErrorBuilder,
	viewer InitViewer,
) InitPresenter {

	return InitPresenter{
		viewableErrorBuilder: viewableErrorBuilder,
		viewer:               viewer,
	}
}

func (i InitPresenter) PresentToView(response features.InitResponse) {
	viewData := views.InitViewData{}

	if response.Error == nil {
		envName := response.Content.EnvName
		envAlreadyCreated := response.Content.EnvAlreadyCreated

		viewDataMessage := "The sandbox \"" + envName + "\" was initialized."
		if envAlreadyCreated {
			viewDataMessage = "The sandbox \"" + envName + "\" is already initialized."
		}

		currentCloudProviderCmd := string(globals.CurrentCloudProvider)
		currentCloudProviderCmdArgs := globals.CurrentCloudProviderArgs

		if len(currentCloudProviderCmdArgs) > 0 {
			currentCloudProviderCmd += " " + currentCloudProviderCmdArgs
		}

		bold := config.ColorsBold

		viewDataSubtext := fmt.Sprintf(
			"The public IP of your sandbox is: %s\n\n"+
				"To connect to your sandbox:\n\n"+
				"  - With your editor: `%s`\n\n"+
				"  - With SSH        : `%s`\n\n"+
				"To allow TCP traffic on a port: `%s`",
			bold(response.Content.EnvPublicIPAddress),
			bold(config.ColorsBlue("eleven "+currentCloudProviderCmd+" edit "+envName)),
			bold(config.ColorsBlue("ssh "+response.Content.EnvLocalSSHConfigHostname)),
			bold(config.ColorsBlue("eleven "+currentCloudProviderCmd+" serve "+envName+" <port> [--as <domain>]")),
		)

		runtimes := response.Content.EnvRuntimes
		runtimesAsText := ""

		for runtimeName, runtimeVersion := range runtimes {
			runtimesAsText += config.ColorsBGBlue(
				fmt.Sprintf(" %s@%s ", runtimeName, runtimeVersion),
			) + " "
		}

		if len(runtimesAsText) > 0 {
			viewDataSubtext += "\n\n"
			viewDataSubtext += "Installed runtimes: " + bold(config.ColorsWhite(runtimesAsText))
		}

		viewData.Content = views.InitViewDataContent{
			ShowAsWarning: envAlreadyCreated,
			Message:       viewDataMessage,
			Subtext:       viewDataSubtext,
		}

		i.viewer.View(viewData)
		return
	}

	viewData.Error = i.viewableErrorBuilder.Build(response.Error)
	i.viewer.View(viewData)
}
