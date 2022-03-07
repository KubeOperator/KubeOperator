package initial

import (
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/phases"
	"github.com/KubeOperator/KubeOperator/pkg/util/kobe"
)

const (
	initNetwork = "09-plugin-network.yml"
)

type NetworkPhase struct {
}

func (NetworkPhase) Name() string {
	return "InitNetwork"
}

func (s NetworkPhase) Run(b kobe.Interface, fileName string) error {
	return phases.RunPlaybookAndGetResult(b, initNetwork, "", fileName)
}
