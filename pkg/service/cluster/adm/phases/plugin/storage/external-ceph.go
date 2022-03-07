package storage

import (
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/phases"
	"github.com/KubeOperator/KubeOperator/pkg/util/kobe"
)

const (
	externalCephStorage = "10-plugin-cluster-storage-external-ceph.yml"
)

type ExternalCephStoragePhase struct {
	ProvisionerName string
}

func (n ExternalCephStoragePhase) Name() string {
	return "CreateExternalCephStorage"
}

func (n ExternalCephStoragePhase) Run(b kobe.Interface, fileName string) error {
	if n.ProvisionerName != "" {
		b.SetVar("storage_rbd_provisioner_name", n.ProvisionerName)
	}
	return phases.RunPlaybookAndGetResult(b, externalCephStorage, "", fileName)
}
