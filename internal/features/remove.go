package features

import (
	"github.com/eleven-sh/eleven/features"
)

type RemoveResponse struct {
	Error   error
	Content RemoveResponseContent
}

type RemoveResponseContent struct {
	EnvName string
}

type RemovePresenter interface {
	PresentToView(RemoveResponse)
}

type RemoveOutputHandler struct {
	presenter RemovePresenter
}

func NewRemoveOutputHandler(
	presenter RemovePresenter,
) RemoveOutputHandler {

	return RemoveOutputHandler{
		presenter: presenter,
	}
}

func (r RemoveOutputHandler) HandleOutput(output features.RemoveOutput) error {
	output.Stepper.StopCurrentStep()

	handleError := func(err error) error {
		r.presenter.PresentToView(RemoveResponse{
			Error: err,
		})

		return err
	}

	if output.Error != nil {
		return handleError(output.Error)
	}

	env := output.Content.Env

	r.presenter.PresentToView(RemoveResponse{
		Content: RemoveResponseContent{
			EnvName: env.Name,
		},
	})

	return nil
}
