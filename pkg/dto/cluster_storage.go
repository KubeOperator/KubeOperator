package dto

import "github.com/KubeOperator/KubeOperator/pkg/model"

type ClusterStorageProvisioner struct {
	model.ClusterStorageProvisioner
	Vars map[string]interface{}
}

type ClusterStorageProvisionerCreation struct {
	Name string
	Type string
	Vars map[string]interface{}
}

type ClusterStorageProvisionerBatch struct {
	Items     []ClusterStorageProvisioner `json:"items"`
	Operation string                      `json:"operation"`
}
