package reset

import (
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/phases"
	"github.com/KubeOperator/KubeOperator/pkg/util/kobe"
)

const (
	resetCluster = "99-reset-cluster.yml"
)

type ResetClusterPhase struct {
}

func (s ResetClusterPhase) Name() string {
	return "ResetCluster"
}

func (s ResetClusterPhase) Run(b kobe.Interface, fileName string) error {
	return phases.RunPlaybookAndGetResult(b, resetCluster, "", fileName)
}
