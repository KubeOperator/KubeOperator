package backup

import (
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/phases"
	"github.com/KubeOperator/KubeOperator/pkg/util/kobe"
)

const (
	backupCluster = "94-backup-cluster.yml"
)

type BackupClusterPhase struct {
}

func (backup BackupClusterPhase) Name() string {
	return "backupCluster"
}

func (backup BackupClusterPhase) Run(b kobe.Interface, fileName string) error {
	return phases.RunPlaybookAndGetResult(b, backupCluster, "", fileName)
}
