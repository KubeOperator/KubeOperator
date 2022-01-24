package dto

import "github.com/KubeOperator/KubeOperator/pkg/model"

type ClusterStorageProvisioner struct {
	model.ClusterStorageProvisioner
	Vars map[string]interface{} `json:"vars"`
}

type ClusterStorageProvisionerCreation struct {
	Name string                 `json:"name" validate:"commonname,required"`
	Type string                 `json:"type" validate:"oneof=nfs external-ceph rook-ceph cinder vsphere glusterfs oceanstor"`
	Vars map[string]interface{} `json:"vars" validate:"-"`
}

type ClusterStorageProvisionerSync struct {
	Name   string `json:"name" validate:"commonname,required"`
	Type   string `json:"type" validate:"oneof=nfs external-ceph rook-ceph cinder vsphere glusterfs oceanstor"`
	Status string `json:"status" validate:"-"`
}

type ClusterStorageProvisionerBatch struct {
	Items     []ClusterStorageProvisioner `json:"items"`
	Operation string                      `json:"operation"`
}
