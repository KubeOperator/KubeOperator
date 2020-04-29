package phases

import (
	"fmt"
	"ko3-gin/pkg/cluster/adm/workflow"
	"ko3-gin/pkg/host"
)

func NewKubeConfigPhase() workflow.Phase {
	return workflow.Phase{
		Name:   "kube-config",
		Phases: nil,
		Run:    runKubeConfig,
	}
}

func runKubeConfig(data workflow.RunData, host host.Host) error {
	fmt.Println("set kube config")
	return nil
}
