package reset

import (
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/phases"
	"github.com/KubeOperator/KubeOperator/pkg/util/kobe"
	"io"
)

const (
	resetCluster = "99-reset-cluster.yml"
)

type ResetClusterPhase struct {
}

func (s ResetClusterPhase) Name() string {
	return "ResetCluster"
}

func (s ResetClusterPhase) Run(b kobe.Interface, writer io.Writer) error {
	return phases.RunPlaybookAndGetResult(b, resetCluster, "", writer)
}
