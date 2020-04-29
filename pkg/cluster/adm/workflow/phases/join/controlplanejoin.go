package phases

import (
	"fmt"
	"ko3-gin/pkg/cluster/adm/workflow"
	"ko3-gin/pkg/host"
)

func NewControlPlaneJoinPhase() workflow.Phase {
	return workflow.Phase{
		Name:   "control-plane-join",
		Phases: nil,
		Run:    runControlPlaneJoin,
	}
}

func runControlPlaneJoin(data workflow.RunData, host host.Host) error {
	fmt.Print("join worker...")
	return nil
}
