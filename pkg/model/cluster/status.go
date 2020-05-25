package cluster

import (
	commonModel "github.com/KubeOperator/KubeOperator/pkg/model/common"
)

type Status struct {
	commonModel.BaseModel
	ID         string
	ClusterID  string
	Version    string
	Message    string
	Phase      string
}

func (s Status) TableName() string {
	return "ko_cluster_status"
}
