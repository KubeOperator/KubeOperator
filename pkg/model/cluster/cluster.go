package cluster

import (
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	uuid "github.com/satori/go.uuid"
	"time"
)

type Cluster struct {
	common.BaseModel
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
