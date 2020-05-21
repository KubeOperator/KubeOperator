package cluster

import (
	commonModel "github.com/KubeOperator/KubeOperator/pkg/model/common"
	uuid "github.com/satori/go.uuid"
	"time"
)

type Condition struct {
	commonModel.BaseModel
	Name          string
	StatusID      string
	Status        string
	Message       string
	LastProbeTime time.Time
}

func (c *Condition) BeforeCreate() error {
	c.ID = uuid.NewV4().String()
	c.CreatedDate = time.Now()
	c.UpdatedDate = time.Now()
	return nil
}

func (c *Condition) BeforeUpdate() error {
	c.UpdatedDate = time.Now()
	return nil
}

func (c Condition) TableName() string {
	return "ko_cluster_condition"
}
