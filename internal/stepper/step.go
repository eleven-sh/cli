package stepper

import (
	"github.com/briandowns/spinner"
	"github.com/eleven-sh/cli/internal/interfaces"
)

type Step struct {
	logger          interfaces.Logger
	spin            *spinner.Spinner
	removeAfterDone bool
}

func (s *Step) Done() {
	s.spin.Stop()

	if !s.removeAfterDone {
		s.logger.Log(s.spin.Prefix + "... done")
	}
}
