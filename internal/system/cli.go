package system

import (
	"bufio"
	"io"
	"strings"

	"github.com/eleven-sh/cli/internal/config"
	"github.com/eleven-sh/cli/internal/interfaces"
)

func AskForConfirmation(
	logger interfaces.Logger,
	stdin io.Reader,
	question string,
) (bool, error) {

	stdinReader := bufio.NewReader(stdin)

	logger.Log(config.ColorsBold(config.ColorsYellow("Warning!") + " " + question))

	logger.Log("\nOnly \"yes\" will be accepted to confirm. (You could use \"--force\" next time).\n")
	logger.LogNoNewline(config.ColorsBold("Confirm? "))

	response, err := stdinReader.ReadString('\n')

	if err != nil {
		return false, err
	}

	sanitizedResponse := strings.TrimSpace(response)

	if sanitizedResponse == "yes" {
		return true, nil
	}

	logger.Log("")

	return false, nil
}
