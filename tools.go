// +build tools

package tools

// See https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module
// for more details

import (
	_ "github.com/go-swagger/go-swagger/cmd/swagger"
	_ "github.com/kyoh86/richgo"
	_ "github.com/myitcv/gobin"
)
