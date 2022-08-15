package dto

import (
	"encoding/json"

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
	Type     string `json:"type"`
	Version  string `json:"version"`
	Describe string `json:"describe"`
	Vars     string `json:"vars"`

	Disabled bool   `json:"disabled"`
	Status   string `json:"status"`
	Message  string `json:"message"`
}

type ComponentCreate struct {
	ClusterName string                 `json:"clusterName"`
	Name        string                 `json:"name"`
	Type        string                 `json:"type"`
	Version     string                 `json:"version"`
	Describe    string                 `json:"describe"`
	Vars        map[string]interface{} `json:"vars"`
}

type ComponentSync struct {
	ClusterName string   `json:"clusterName"`
	Names       []string `json:"names"`
}

func (c ComponentCreate) ComponentCreate2Mo() model.ClusterSpecComponent {
	vars, _ := json.Marshal(c.Vars)
	return model.ClusterSpecComponent{
		Name:    c.Name,
		Type:    c.Type,
		Version: c.Version,
		Status:  constant.StatusCreating,
		Vars:    string(vars),
	}
}
