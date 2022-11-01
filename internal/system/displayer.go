package system

import (
	"fmt"
	"io"
)

type Displayer struct{}

func NewDisplayer() Displayer {
	return Displayer{}
}

func (Displayer) Display(w io.Writer, format string, args ...interface{}) {
	toDisplay := replaceNewLinesForOS(
		fmt.Sprintf(format, args...),
	)

	fmt.Fprint(w, toDisplay)
}
