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
	fmt.Fprintf(os.Stderr, config.ColorsCyan(format)+"\n", v...)
}

func (Logger) Warning(format string, v ...interface{}) {
	fmt.Fprintf(os.Stderr, config.ColorsYellow(format)+"\n", v...)
}

func (Logger) Error(format string, v ...interface{}) {
	fmt.Fprintf(os.Stderr, config.ColorsRed(format)+"\n", v...)
}

func (Logger) Log(format string, v ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", v...)
}

func (Logger) LogNoNewline(format string, v ...interface{}) {
	fmt.Fprintf(os.Stderr, format, v...)
}

func (l Logger) Write(p []byte) (n int, err error) {
	l.Log(string(p))
	return len(p), nil
}
