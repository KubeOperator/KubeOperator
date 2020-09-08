package backup

import (
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/phases"
	"github.com/KubeOperator/KubeOperator/pkg/util/kobe"
)

const (
	backupCluster = "95-restore-cluster-custom.yml"
)

type BackupClusterPhase struct {
	backupFileName string
}

func (backup BackupClusterPhase) Name() string {
	return "backupCluster"
}

func (backup BackupClusterPhase) Run(b kobe.Interface) error {
	if backup.backupFileName != "" {
		b.SetVar("custom_etcd_snapshot", backup.backupFileName)
	}
	return phases.RunPlaybookAndGetResult(b, backupCluster)
}
