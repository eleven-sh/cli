package system

import (
	"fmt"
	"os"

	"github.com/eleven-sh/cli/internal/config"
)

type Logger struct{}

func NewLogger() Logger {
	return Logger{}
}

func (Logger) Info(format string, v ...interface{}) {
	toDisplay := replaceNewLinesForOS(
		fmt.Sprintf(format+"\n", v...),
	)

	fmt.Fprint(os.Stderr, config.ColorsCyan(toDisplay))
}

func (Logger) Warning(format string, v ...interface{}) {
	toDisplay := replaceNewLinesForOS(
		fmt.Sprintf(format+"\n", v...),
	)

	fmt.Fprint(os.Stderr, config.ColorsYellow(toDisplay))
}

func (Logger) Error(format string, v ...interface{}) {
	toDisplay := replaceNewLinesForOS(
		fmt.Sprintf(format+"\n", v...),
	)

	fmt.Fprint(os.Stderr, config.ColorsRed(toDisplay))
}

func (Logger) Log(format string, v ...interface{}) {
	toDisplay := replaceNewLinesForOS(
		fmt.Sprintf(format+"\n", v...),
	)

	fmt.Fprint(os.Stderr, toDisplay)
}

func (Logger) LogNoNewline(format string, v ...interface{}) {
	toDisplay := replaceNewLinesForOS(
		fmt.Sprintf(format, v...),
	)

	fmt.Fprint(os.Stderr, toDisplay)
}

func (l Logger) Write(p []byte) (n int, err error) {
	l.Log(string(p))
	return len(p), nil
}
