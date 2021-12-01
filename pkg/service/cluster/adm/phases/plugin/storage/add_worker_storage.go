package storage

import (
	"io"

	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/phases"
	"github.com/KubeOperator/KubeOperator/pkg/util/kobe"
)

const (
	AddWorkerStorage = "10-plugin-cluster-storage-add-worker.yml"
)

type AddWorkerStoragePhase struct {
	EnableNfsProvisioner          string
	NfsVersion                    string
	EnableGfsProvisioner          string
	EnableExternalCephProvisioner string
	AddWorker                     bool
}

func (n AddWorkerStoragePhase) Name() string {
	return "CreateNfsStorage"
}

func (n AddWorkerStoragePhase) Run(b kobe.Interface, writer io.Writer) error {
	b.SetVar("enable_nfs_provisioner", n.EnableNfsProvisioner)
	if n.EnableNfsProvisioner == "disable" {
		b.SetVar("storage_nfs_server_version", n.NfsVersion)
	}
	b.SetVar("enable_gfs_provisioner", n.EnableGfsProvisioner)
	b.SetVar("enable_external_ceph_provisioner", n.EnableExternalCephProvisioner)
	var tag string
	if n.AddWorker {
		tag = "add_worker"
	}

	return phases.RunPlaybookAndGetResult(b, AddWorkerStorage, tag, writer)
}
