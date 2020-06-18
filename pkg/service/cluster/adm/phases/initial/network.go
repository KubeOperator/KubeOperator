package initial

import (
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/phases"
	"github.com/KubeOperator/KubeOperator/pkg/util/kobe"
)

const (
	initNetwork = "10-network-plugin.yml"
)

type NetworkPhase struct {
}

func (NetworkPhase) Name() string {
	return "InitNetwork"
}

func (s NetworkPhase) Run(b kobe.Interface) error {
	return phases.RunPlaybookAndGetResult(b, initNetwork)
}
