package dto

import "github.com/KubeOperator/KubeOperator/pkg/model"

type ClusterTool struct {
	model.ClusterTool
	Vars map[string]interface{} `json:"vars"`
}
