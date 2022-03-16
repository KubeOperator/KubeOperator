package model

import (
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	uuid "github.com/satori/go.uuid"
)

type ClusterGpu struct {
	common.BaseModel
	ID        string `json:"id" gorm:"type:varchar(64)"`
	ClusterID string `json:"cluster_id"`
	Describe  string `json:"describe"`
	Status    string `json:"status"`
	Message   string `json:"message" gorm:"type:text(65535)"`
	Vars      string `json:"vars" gorm:"type:text(65535)"`
}

func (c *ClusterGpu) BeforeCreate() (err error) {
	c.ID = uuid.NewV4().String()
	return nil
}
