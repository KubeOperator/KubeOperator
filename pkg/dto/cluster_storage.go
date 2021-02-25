package dto

import "github.com/KubeOperator/KubeOperator/pkg/model"

type ClusterStorageProvisioner struct {
	model.ClusterStorageProvisioner
	Vars map[string]interface{} `json:"vars"`
}

type ClusterStorageProvisionerCreation struct {
	Name string                 `json:"name"`
	Type string                 `json:"type"`
	Vars map[string]interface{} `json:"vars"`
}

type ClusterStorageProvisionerSync struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	Status string `json:"status"`
}

type ClusterStorageProvisionerBatch struct {
	Items     []ClusterStorageProvisioner `json:"items"`
	Operation string                      `json:"operation"`
}
