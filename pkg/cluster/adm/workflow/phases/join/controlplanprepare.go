package phases

import (
	"fmt"
	"ko3-gin/pkg/cluster/adm/workflow"
	"ko3-gin/pkg/host"
)

func NewControlPlanePreparePhase() workflow.Phase {
	return workflow.Phase{
		Name:   "control-plane-join",
		Phases: nil,
		Run:    runControlPlaneJoin,
	}
}

func runControlPlanePrepare(data workflow.RunData, host host.Host) error {
	fmt.Println("prepare control plane")
	return nil
}
