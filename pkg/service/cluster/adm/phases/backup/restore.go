package backup

import (
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/phases"
	"github.com/KubeOperator/KubeOperator/pkg/util/kobe"
)

const (
	restoreCluster = "95-restore-cluster.yml"
)

type RestoreClusterPhase struct {
}

func (restore RestoreClusterPhase) Name() string {
	return "backupCluster"
}

func (restore RestoreClusterPhase) Run(b kobe.Interface, fileName string) error {
	return phases.RunPlaybookAndGetResult(b, restoreCluster, "", fileName)
}
