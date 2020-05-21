package cluster

import (
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
)

type Spec struct {
	common.BaseModel
	ID          string
	ClusterID   string
	Version     string
	NetworkType string
	ClusterCIDR string `gorm:"column:cluster_cidr"`
	ServiceCIDR string `gorm:"column:service_cidr"`
	Nodes       []Node
}

func (s Spec) TableName() string {
	return "ko_cluster_spec"
}
