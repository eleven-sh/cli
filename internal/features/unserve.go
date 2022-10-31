package features

import (
	"github.com/eleven-sh/cli/internal/agent"
	"github.com/eleven-sh/cli/internal/entities"
	"github.com/eleven-sh/eleven/features"
)

type UnserveResponse struct {
	Error   error
	Content UnserveResponseContent
}

type UnserveResponseContent struct {
	EnvName string
	Port    string
}

type UnservePresenter interface {
	PresentToView(UnserveResponse)
}

type UnserveOutputHandler struct {
	presenter          UnservePresenter
	agentClientBuilder agent.ClientBuilder
}

func NewUnserveOutputHandler(
	presenter UnservePresenter,
	agentClientBuilder agent.ClientBuilder,
) UnserveOutputHandler {

	return UnserveOutputHandler{
		presenter:          presenter,
		agentClientBuilder: agentClientBuilder,
	}
}

func (u UnserveOutputHandler) HandleOutput(output features.UnserveOutput) error {
	stepper := output.Stepper

	handleError := func(err error) error {
		stepper.StopCurrentStep()

		u.presenter.PresentToView(UnserveResponse{
			Error: err,
		})

		return err
	}

	if output.Error != nil {
		return handleError(output.Error)
	}

	env := output.Content.Env

	err := entities.ReconcileServedPortsState(
		env,
		u.agentClientBuilder,
	)

	if err != nil {
		return handleError(err)
	}

	stepper.StopCurrentStep()

	u.presenter.PresentToView(UnserveResponse{
		Content: UnserveResponseContent{
			EnvName: env.Name,
			Port:    output.Content.Port,
		},
	})

	return nil
}
