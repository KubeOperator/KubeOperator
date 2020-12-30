package storage

import (
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/phases"
	"github.com/KubeOperator/KubeOperator/pkg/util/kobe"
	"io"
)

const (
	NfsStorage = "10-plugin-cluster-storage-nfs.yml"
)

type NfsStoragePhase struct {
	NfsServerVersion string
	NfsServer        string
	NfsServerPath    string
	ProvisionerName  string
}

func (n NfsStoragePhase) Name() string {
	return "CrateNfsStorage"
}

func (n NfsStoragePhase) Run(b kobe.Interface, writer io.Writer) error {
	if n.NfsServerVersion != "" {
		b.SetVar("storage_nfs_server_version", n.NfsServerVersion)
	}
	if n.NfsServer != "" {
		b.SetVar("storage_nfs_server", n.NfsServer)
	}
	if n.NfsServerPath != "" {
		b.SetVar("storage_nfs_server_path", n.NfsServerPath)
	}
	if n.ProvisionerName != "" {
		b.SetVar("storage_nfs_provisioner_name", n.ProvisionerName)
	}
	return phases.RunPlaybookAndGetResult(b, NfsStorage, "", writer)
}
