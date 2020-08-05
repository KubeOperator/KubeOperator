package storage

import (
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/phases"
	"github.com/KubeOperator/KubeOperator/pkg/util/kobe"
)

const (
	externalCephStorage = "10-plugin-cluster-storage-external-ceph.yml"
)

type ExternalCephStoragePhase struct {
	CephMonitor               string
	CephOsdPool               string
	CephAdminId               string
	CephAdminSecret           string
	CephUserId                string
	CephUserSecret            string
	CephFsType                string
	CephImageFormat           string
	StorageRbdProvisionerName string
	ProvisionerName           string
}

func (n ExternalCephStoragePhase) Name() string {
	return "CreateExternalCephStorage"
}

func (n ExternalCephStoragePhase) Run(b kobe.Interface) error {
	if n.CephMonitor != "" {
		b.SetVar("ceph_monitor", n.CephMonitor)
	}
	if n.CephOsdPool != "" {
		b.SetVar("ceph_osd_pool", n.CephOsdPool)
	}
	if n.CephAdminId != "" {
		b.SetVar("ceph_admin_id", n.CephAdminId)
	}
	if n.CephAdminSecret != "" {
		b.SetVar("ceph_admin_secret", n.CephAdminSecret)
	}
	if n.CephUserId != "" {
		b.SetVar("ceph_user_id", n.CephUserId)
	}
	if n.CephUserSecret != "" {
		b.SetVar("ceph_user_secret", n.CephUserSecret)
	}
	if n.CephFsType != "" {
		b.SetVar("ceph_fsType", n.CephFsType)
	}
	if n.CephImageFormat != "" {
		b.SetVar("ceph_imageFormat", n.CephImageFormat)
	}
	if n.ProvisionerName != "" {
		b.SetVar("storage_rbd_provisioner_name", n.ProvisionerName)
	}
	return phases.RunPlaybookAndGetResult(b, externalCephStorage)
}
