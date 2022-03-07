package storage

import (
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/phases"
	"github.com/KubeOperator/KubeOperator/pkg/util/kobe"
)

const (
	rookCephStorage = "10-plugin-cluster-storage-rook-ceph.yml"
)

type RookCephStoragePhase struct {
	StorageRookPath string
}

func (n RookCephStoragePhase) Name() string {
	return "CreateRookCephStorage"
}

func (n RookCephStoragePhase) Run(b kobe.Interface, fileName string) error {
	if n.StorageRookPath != "" {
		b.SetVar("storage_rook_path", n.StorageRookPath)
	}
	b.SetVar("storage_rook_enabled", "true")
	return phases.RunPlaybookAndGetResult(b, rookCephStorage, "", fileName)
}
