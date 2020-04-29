package workflow

import (
	"ko3-gin/pkg/host"
)

type Phase struct {
	Name   string
	Phases []Phase
	Run    func(data RunData, host host.Host) error
}
