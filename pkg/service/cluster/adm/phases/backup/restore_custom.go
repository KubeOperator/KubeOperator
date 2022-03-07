package backup

import (
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/phases"
	"github.com/KubeOperator/KubeOperator/pkg/util/kobe"
)

const (
	restoreClusterCustom = "95-restore-cluster-custom.yml"
)

type RestoreClusterCustomPhase struct {
	BackupFileName string
}

func (restore RestoreClusterCustomPhase) Name() string {
	return "backupCluster"
}

func (restore RestoreClusterCustomPhase) Run(b kobe.Interface, fileName string) error {
	if restore.BackupFileName != "" {
		b.SetVar("custom_etcd_snapshot", restore.BackupFileName)
	}
	return phases.RunPlaybookAndGetResult(b, restoreClusterCustom, "", fileName)
}
