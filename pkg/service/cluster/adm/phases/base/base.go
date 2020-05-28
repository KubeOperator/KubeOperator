package phase

import (
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/phases"
	"github.com/KubeOperator/KubeOperator/pkg/util/kobe"
)

const (
	playbookNameBase = "01-base.yml"
)

type SystemConfigPhase struct {
}

func (s SystemConfigPhase) Name() string {
	return "ConfigSystem"
}

func (s SystemConfigPhase) Run(b kobe.Interface) (result kobe.Result, err error) {
	return phases.RunPlaybookAndGetResult(b, playbookNameBase)
}
