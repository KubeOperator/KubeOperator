package dto

import (
	"encoding/json"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/model"
)

type ClusterStorageProvisioner struct {
	model.ClusterStorageProvisioner
	Vars map[string]interface{} `json:"vars"`
}

type ClusterStorageProvisionerCreation struct {
	Name      string                 `json:"name"`
	Namespace string                 `json:"namespace"`
	Type      string                 `json:"type"`
	Vars      map[string]interface{} `json:"vars"`

	IsInCluster bool `json:"isInCluster"`
}

type ClusterStorageProvisionerSync struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Type      string `json:"type"`
	Status    string `json:"status"`
}

type ClusterStorageProvisionerBatch struct {
	Items     []ClusterStorageProvisioner `json:"items"`
	Operation string                      `json:"operation"`
}

type ClusterStorageProvisionerLoad struct {
	Name   string                 `json:"name"`
	Type   string                 `json:"type"`
	Status string                 `json:"status"`
	Vars   map[string]interface{} `json:"vars"`
}

type DeploymentSearch struct {
	ApiServer string `json:"apiServer"`
	Router    string `json:"router"`
	Token     string `json:"token"`
	Namespace string `json:"namespace"`
}

func (c ClusterStorageProvisionerCreation) ProvisionerCreate2Mo() model.ClusterStorageProvisioner {
	vars, _ := json.Marshal(c.Vars)
	provisioner := model.ClusterStorageProvisioner{
		Name:      c.Name,
		Namespace: c.Namespace,
		Type:      c.Type,
		Vars:      string(vars),
		Status:    constant.StatusCreating,
	}
	if c.IsInCluster {
		provisioner.Status = constant.StatusRunning
	}
	return provisioner
}
