package dto

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/model"
)

type ComponentPage struct {
	Items []Component `json:"items"`
	Total int         `json:"total"`
}

type Component struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Version  string `json:"version"`
	Describe string `json:"describe"`

	Status  string `json:"status"`
	Message string `json:"message"`
}

type ComponentCreate struct {
	ClusterName string `json:"clusterName"`
	Name        string `json:"name"`
	Version     string `json:"version"`
	Describe    string `json:"describe"`
}

type ComponentDelete struct {
	ClusterName string `json:"clusterName"`
	Name        string `json:"name"`
	Version     string `json:"version"`
	Describe    string `json:"describe"`
}

func (c ComponentCreate) ComponentCreate2Mo() model.ClusterSpecComponent {
	return model.ClusterSpecComponent{
		Name:    c.Name,
		Version: c.Version,
		Status:  constant.StatusDisabled,
	}
}
