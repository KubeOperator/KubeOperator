package backup

import (
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/phases"
	"github.com/KubeOperator/KubeOperator/pkg/util/kobe"
	"io"
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

func (restore RestoreClusterCustomPhase) Run(b kobe.Interface, writer io.Writer) error {
	if restore.BackupFileName != "" {
		b.SetVar("custom_etcd_snapshot", restore.BackupFileName)
	}
	return phases.RunPlaybookAndGetResult(b, restoreClusterCustom, "", writer)
}
