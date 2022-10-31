package features

import (
	"os"

	"github.com/eleven-sh/cli/internal/system"
	"github.com/eleven-sh/eleven/features"
)

type UninstallResponse struct {
	Error   error
	Content UninstallResponseContent
}

type UninstallResponseContent struct {
	ElevenAlreadyUninstalled  bool
	SuccessMessage            string
	AlreadyUninstalledMessage string
	ElevenExecutablePath      string
	ElevenConfigDirPath       string
}

type UninstallPresenter interface {
	PresentToView(UninstallResponse)
}

type UninstallOutputHandler struct {
	presenter UninstallPresenter
}

func NewUninstallOutputHandler(
	presenter UninstallPresenter,
) UninstallOutputHandler {

	return UninstallOutputHandler{
		presenter: presenter,
	}
}

func (u UninstallOutputHandler) HandleOutput(output features.UninstallOutput) error {
	output.Stepper.StopCurrentStep()

	handleError := func(err error) error {
		u.presenter.PresentToView(UninstallResponse{
			Error: err,
		})

		return err
	}

	if output.Error != nil {
		return handleError(output.Error)
	}

	elevenExecutablePath, err := os.Executable()

	if err != nil {
		elevenExecutablePath = "<path_not_found>"
	}

	elevenConfigDirPath := system.UserConfigDir()

	u.presenter.PresentToView(UninstallResponse{
		Content: UninstallResponseContent{
			ElevenAlreadyUninstalled:  output.Content.ElevenAlreadyUninstalled,
			SuccessMessage:            output.Content.SuccessMessage,
			AlreadyUninstalledMessage: output.Content.AlreadyUninstalledMessage,
			ElevenExecutablePath:      elevenExecutablePath,
			ElevenConfigDirPath:       elevenConfigDirPath,
		},
	})

	return nil
}
