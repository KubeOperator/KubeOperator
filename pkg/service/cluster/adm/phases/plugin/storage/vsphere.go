package storage

import (
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/phases"
	"github.com/KubeOperator/KubeOperator/pkg/util/kobe"
)

const (
	vsphereStorage = "10-plugin-cluster-storage-vsphere.yml"
)

type VsphereStoragePhase struct {
	VcUsername        string
	VcPassword        string
	VcHost            string
	VcPort            string
	Datacenter        string
	Folder            string
	ProvisionerName   string
	StoragePolicyName string
}

func (n VsphereStoragePhase) Name() string {
	return "CreateVsphereStorage"
}

func (n VsphereStoragePhase) Run(b kobe.Interface, fileName string) error {
	if n.VcUsername != "" {
		b.SetVar("vc_username", n.VcUsername)
	}

	if n.VcPassword != "" {
		b.SetVar("vc_password", n.VcPassword)
	}

	if n.VcHost != "" {
		b.SetVar("vc_host", n.VcHost)
	}

	if n.VcPort != "" {
		b.SetVar("vc_port", n.VcPort)
	}

	if n.Datacenter != "" {
		b.SetVar("datacenter", n.Datacenter)
	}

	if n.Folder != "" {
		b.SetVar("folder", n.Folder)
	}

	if n.StoragePolicyName != "" {
		b.SetVar("vc_storage_policy", n.StoragePolicyName)
	}

	return phases.RunPlaybookAndGetResult(b, vsphereStorage, "", fileName)
}
