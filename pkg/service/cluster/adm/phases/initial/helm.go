package initial

import (
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/phases"
	"github.com/KubeOperator/KubeOperator/pkg/util/kobe"
)

const (
	initHelm = "11-helm-install.yml"
)

type HelmPhase struct {
}

func (h HelmPhase) Name() string {
	return "InitHelm"
}

func (h HelmPhase) Run(b kobe.Interface, fileName string) error {
	return phases.RunPlaybookAndGetResult(b, initHelm, "", fileName)
}
