package cluster

import (
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	"time"
)

type Condition struct {
	common.BaseModel
	ID            string
	Name          string
	StatusID      string
	Status        string
	Message       string
	LastProbeTime time.Time
}

func (c Condition) TableName() string {
	return "ko_cluster_condition"
}
