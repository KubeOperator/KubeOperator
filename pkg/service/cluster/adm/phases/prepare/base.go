package prepare

import (
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/phases"
	"github.com/KubeOperator/KubeOperator/pkg/util/kobe"
)

const (
	prepareBase = "01-base.yml"
)

type BaseSystemConfigPhase struct {
}

func (s BaseSystemConfigPhase) Name() string {
	return "BasicConfigSystem"
}

func (s BaseSystemConfigPhase) Run(b kobe.Interface) (err error) {
	return phases.RunPlaybookAndGetResult(b, prepareBase)
}
