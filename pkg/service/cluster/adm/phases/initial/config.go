package initial

import (
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/phases"
	"github.com/KubeOperator/KubeOperator/pkg/util/kobe"
)

const (
	initKubeConfig = "07-kubernetes-config.yml"
)

type KubeConfigPhase struct {
}

func (KubeConfigPhase) Name() string {
	return "InitKubernetesConfig"
}

func (s KubeConfigPhase) Run(b kobe.Interface) (result kobe.Result, err error) {
	return phases.RunPlaybookAndGetResult(b, initKubeConfig)
}
