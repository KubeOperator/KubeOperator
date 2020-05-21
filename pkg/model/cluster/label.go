package cluster

import (
	commonModel "github.com/KubeOperator/KubeOperator/pkg/model/common"
)

type Label struct {
	commonModel.BaseModel
	ID   string
	Name   string
	Value  string
	NodeID string
}

func (l Label) TableName() string {
	return "ko_cluster_node_label"
}
