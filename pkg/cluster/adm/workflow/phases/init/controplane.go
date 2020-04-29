package phases

import (
	"ko3-gin/pkg/cluster/adm/workflow"
	"ko3-gin/pkg/host"
)

func NewControlPlanePhase() workflow.Phase {
	return workflow.Phase{
		Name:   "control-plane",
		Phases: nil,
		Run:    runControlPlane,
	}
}

func runControlPlane(data workflow.RunData, host host.Host) error {
	return nil
}
