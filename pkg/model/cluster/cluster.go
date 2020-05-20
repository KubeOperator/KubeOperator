package cluster

import (
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	uuid "github.com/satori/go.uuid"
	"time"
)

type Spec struct {
	Version     string
	NetworkType string
	ClusterCIDR string
	ServiceCIDR string
}



type Cluster struct {
	common.BaseModel
	Spec
	Nodes []Node
}

func (c *Cluster) BeforeCreate() error {
	c.ID = uuid.NewV4().String()
	c.CreatedDate = time.Now()
	c.UpdatedDate = time.Now()
	return nil
}

func (c *Cluster) BeforeUpdate() error {
	c.UpdatedDate = time.Now()
	return nil
}

func (c Cluster) TableName() string {
	return "ko_cluster"
}
