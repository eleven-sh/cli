package features

import (
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/eleven-sh/cli/internal/agent"
	cliEntities "github.com/eleven-sh/cli/internal/entities"
	"github.com/eleven-sh/cli/internal/hooks"
	"github.com/eleven-sh/eleven/entities"
	"github.com/eleven-sh/eleven/features"
)

type ServeResponse struct {
	Error   error
	Content ServeResponseContent
}

type ServeResponseContent struct {
	EnvName            string
	EnvPublicIPAddress string
	Port               string
	PortBinding        string
}

type ServePresenter interface {
	PresentToView(ServeResponse)
}

type ServeOutputHandler struct {
	presenter          ServePresenter
	agentClientBuilder agent.ClientBuilder
}

func NewServeOutputHandler(
	presenter ServePresenter,
	agentClientBuilder agent.ClientBuilder,
) ServeOutputHandler {

	return ServeOutputHandler{
		presenter:          presenter,
		agentClientBuilder: agentClientBuilder,
	}
}

func (s ServeOutputHandler) HandleOutput(output features.ServeOutput) error {
	stepper := output.Stepper

	handleError := func(err error) error {
		stepper.StopCurrentStep()

		s.presenter.PresentToView(ServeResponse{
			Error: err,
		})

		return err
	}

	if output.Error != nil {
		return handleError(output.Error)
	}

	env := output.Content.Env

	err := cliEntities.ReconcileServedPortsState(
		env,
		s.agentClientBuilder,
	)

	if err != nil {
		return handleError(err)
	}

	portBinding := output.Content.PortBinding

	if !govalidator.IsPort(portBinding) {
		stepper.StartTemporaryStep("Waiting for Let's Encrypt to issue certificate")

		err := hooks.WaitUntilDomainIsReachableViaHTTPS(
			portBinding,
			5*time.Minute,
		)

		if err != nil {
			return handleError(entities.ErrLetsEncryptTimedOut{
				Domain:        portBinding,
				ReturnedError: err,
			})
		}
	}

	stepper.StopCurrentStep()

	s.presenter.PresentToView(ServeResponse{
		Content: ServeResponseContent{
			EnvName:            env.Name,
			EnvPublicIPAddress: env.InstancePublicIPAddress,
			Port:               output.Content.Port,
			PortBinding:        portBinding,
		},
	})

	return nil
}
