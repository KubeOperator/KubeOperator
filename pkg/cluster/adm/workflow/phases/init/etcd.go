package phases

import (
	"ko3-gin/pkg/cluster/adm/workflow"
	"ko3-gin/pkg/host"
)

func NewEtcdPhase() workflow.Phase {
	return workflow.Phase{
		Name:   "etcd",
		Phases: nil,
		Run:    runEtcd,
	}
}

func runEtcd(data workflow.RunData, host host.Host) error {
	return nil
}
