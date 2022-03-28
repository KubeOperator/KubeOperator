package dto

import "github.com/KubeOperator/KubeOperator/pkg/model"

type ClusterTool struct {
	model.ClusterTool
	Vars map[string]interface{} `json:"vars"`
}

type ToolPort struct {
	NodeHost string `json:"nodeHost"`
	NodePort int32  `json:"nodePort"`
}
