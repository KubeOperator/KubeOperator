package cluster

import (
	commonModel "github.com/KubeOperator/KubeOperator/pkg/model/common"
	uuid "github.com/satori/go.uuid"
	"time"
)

type Status struct {
	commonModel.BaseModel
	Version    string
	Message    string
	Phase      string
	Conditions []Condition
}

func (n *Status) BeforeCreate() error {
	n.ID = uuid.NewV4().String()
	n.CreatedDate = time.Now()
	n.UpdatedDate = time.Now()
	return nil
}

func (n *Status) BeforeUpdate() error {
	n.UpdatedDate = time.Now()
	return nil
}

func (n Status) TableName() string {
	return "ko_cluster_status"
}
