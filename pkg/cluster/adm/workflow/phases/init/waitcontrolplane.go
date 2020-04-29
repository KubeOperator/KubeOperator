package phases

import (
	"ko3-gin/pkg/cluster/adm/workflow"
	"ko3-gin/pkg/host"
)

func NewWaitControlPlanePhase() workflow.Phase {
	return workflow.Phase{
		Name:   "wait-control-plane",
		Phases: nil,
		Run:    runWaitControlPlane,
	}
}

func runWaitControlPlane(data workflow.RunData, host host.Host) error {
	return nil
}
