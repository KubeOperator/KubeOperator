package initial

import (
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/phases"
	"github.com/KubeOperator/KubeOperator/pkg/util/kobe"
)

const (
	initNpd = "12-npd.yml"
)

type NpdPhase struct {
}

func (s NpdPhase) Name() string {
	return "Npd Init"
}

func (s NpdPhase) Run(b kobe.Interface, fileName string) error {
	return phases.RunPlaybookAndGetResult(b, initNpd, "", fileName)
}
