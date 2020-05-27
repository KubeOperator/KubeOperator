package cluster

import (
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	uuid "github.com/satori/go.uuid"
	"time"
)

type Condition struct {
	common.BaseModel
	ID            string
	Name          string
	StatusID      string
	Status        string
	Message       string
	OrderNum      int
	LastProbeTime time.Time
}

func (c *Condition) BeforeCreate() (err error) {
	c.ID = uuid.NewV4().String()
	return nil
}

func (c Condition) TableName() string {
	return "ko_cluster_condition"
}
