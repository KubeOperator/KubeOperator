package npd

import (
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/phases"
	"github.com/KubeOperator/KubeOperator/pkg/util/kobe"
)

const (
	npdPlaybook = "12-npd.yml"
)

type NpdPhase struct {
}

func (NpdPhase) Name() string {
	return "Npd"
}

func (c NpdPhase) Run(b kobe.Interface, fileName string) error {
	return phases.RunPlaybookAndGetResult(b, npdPlaybook, "", fileName)
}
