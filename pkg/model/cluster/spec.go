package cluster

import (
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	uuid "github.com/satori/go.uuid"
	"time"
)

type Spec struct {
	common.BaseModel
	Version     string
	NetworkType string
	ClusterCIDR string
	ServiceCIDR string
	Nodes       []Node
}

func (n *Spec) BeforeCreate() error {
	n.ID = uuid.NewV4().String()
	n.CreatedDate = time.Now()
	n.UpdatedDate = time.Now()
	return nil
}

func (n *Spec) BeforeUpdate() error {
	n.UpdatedDate = time.Now()
	return nil
}

func (n Spec) TableName() string {
	return "ko_cluster_spec"
}
