package views

import (
	"io"
	"os"
	"strings"

	"github.com/eleven-sh/cli/internal/config"
)

//go:generate go run github.com/golang/mock/mockgen -destination=../mocks/views_displayer.go -package=mocks github.com/eleven-sh/cli/internal/views Displayer
type Displayer interface {
	Display(w io.Writer, format string, args ...interface{})
}

type BaseView struct {
	Displayer Displayer
}

func NewBaseView(displayer Displayer) BaseView {
	return BaseView{
		Displayer: displayer,
	}
}

func (b BaseView) showErrorView(
	err *ViewableError,
	startWithNewLine bool,
) {

	bold := config.ColorsBold
	red := config.ColorsRed

	if startWithNewLine {
		b.Displayer.Display(
			os.Stdout,
			"\n",
		)
	}

	if len(err.Logs) > 0 {
		b.Displayer.Display(
			os.Stdout,
			"%s\n\n",
			strings.TrimSuffix(err.Logs, "\n"),
		)
	}

	b.Displayer.Display(
		os.Stdout,
		"%s %s\n\n%s\n\n",
		bold(red("Error!")),
		bold(err.Title),
		err.Message,
	)
}

func (b BaseView) ShowErrorView(err *ViewableError) {
	b.showErrorView(err, false)
}

func (b BaseView) ShowErrorViewWithStartingNewLine(err *ViewableError) {
	b.showErrorView(err, true)
}

func (b BaseView) ShowWarningView(warningText, subtext string) {
	bold := config.ColorsBold
	yellow := config.ColorsYellow

	if len(subtext) > 0 {
		b.Displayer.Display(
			os.Stdout,
			"%s %s\n\n%s\n\n",
			bold(yellow("Warning!")),
			bold(warningText),
			subtext,
		)

		return
	}

	b.Displayer.Display(
		os.Stdout,
		"%s %s\n\n",
		bold(yellow("Warning!")),
		bold(warningText),
	)
}

func (b BaseView) ShowSuccessView(successText, subtext string) {
	bold := config.ColorsBold
	green := config.ColorsGreen

	if len(subtext) > 0 {
		b.Displayer.Display(
			os.Stdout,
			"%s %s\n\n%s\n\n",
			bold(green("Success!")),
			bold(successText),
			subtext,
		)

		return
	}

	b.Displayer.Display(
		os.Stdout,
		"%s %s\n\n",
		bold(green("Success!")),
		bold(successText),
	)
}
