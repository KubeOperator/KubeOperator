// +build tools

// This file just exists to ensure we download the tools we need for building
// See https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module

package helm_operator

import (
	_ "github.com/fluxcd/helm-operator/pkg/install"
)
