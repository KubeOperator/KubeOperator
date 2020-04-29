package phases

import (
	"fmt"
	"ko3-gin/pkg/cluster/adm/workflow"
	"ko3-gin/pkg/host"
)

func NewCertsPhase() workflow.Phase {
	return workflow.Phase{
		Name:   "certs",
		Phases: nil,
		Run:    runCerts,
	}
}

func runCerts(data workflow.RunData, host host.Host) error {
	fmt.Println("generate certs...")
	return nil
}
