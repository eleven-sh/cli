package stepper

import (
	"time"

	"github.com/briandowns/spinner"
	"github.com/eleven-sh/cli/internal/config"
	"github.com/eleven-sh/cli/internal/interfaces"
	"github.com/eleven-sh/eleven/stepper"
)

var currentStep *Step

type Stepper struct {
	logger interfaces.Logger
}

func NewStepper(
	logger interfaces.Logger,
) Stepper {

	return Stepper{
		logger: logger,
	}
}

func (s Stepper) startStep(
	step string,
	removeAfterDone bool,
	noNewLineAtStart bool,
) stepper.Step {

	if currentStep == nil && !noNewLineAtStart {
		s.logger.Log("")
	}

	if currentStep != nil {
		currentStep.Done()
		currentStep = nil
	}

	bold := config.ColorsBold

	spin := spinner.New(spinner.CharSets[26], 400*time.Millisecond)
	spin.Prefix = bold(step)
	spin.Start()

	currentStep = &Step{
		logger:          s.logger,
		spin:            spin,
		removeAfterDone: removeAfterDone,
	}

	return currentStep
}

func (s Stepper) StartStep(
	step string,
) stepper.Step {

	removeAfterDone := false
	noNewLineAtStart := false

	return s.startStep(
		step,
		removeAfterDone,
		noNewLineAtStart,
	)
}

func (s Stepper) StartTemporaryStep(
	step string,
) stepper.Step {

	removeAfterDone := true
	noNewLineAtStart := false

	return s.startStep(
		step,
		removeAfterDone,
		noNewLineAtStart,
	)
}

func (s Stepper) StartTemporaryStepWithoutNewLine(
	step string,
) stepper.Step {

	removeAfterDone := true
	noNewLineAtStart := true

	return s.startStep(
		step,
		removeAfterDone,
		noNewLineAtStart,
	)
}

func (s Stepper) StopCurrentStep() {

	if currentStep != nil {
		currentStep.Done()
		currentStep = nil
	}
}
