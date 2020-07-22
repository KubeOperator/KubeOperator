package initial

import (
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/phases"
	"github.com/KubeOperator/KubeOperator/pkg/util/kobe"
)

const (
	initMaster = "07-kubernetes-master.yml"
)

type MasterPhase struct {
}

func (MasterPhase) Name() string {
	return "InitEtcd"
}

func (s MasterPhase) Run(b kobe.Interface) error {
	return phases.RunPlaybookAndGetResult(b, initMaster)
}
