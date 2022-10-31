//go:build tools
// +build tools

// following https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module

package internal

import (
	_ "github.com/golang/mock/mockgen"
	_ "github.com/google/wire/cmd/wire"
)
