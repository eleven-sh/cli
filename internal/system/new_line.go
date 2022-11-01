package system

import (
	"runtime"
	"strings"
)

var NewLineChar = "\n"

func init() {
	if runtime.GOOS == "windows" {
		NewLineChar = "\r\n"
	}
}

func replaceNewLinesForOS(s string) string {
	return strings.ReplaceAll(s, "\n", NewLineChar)
}
