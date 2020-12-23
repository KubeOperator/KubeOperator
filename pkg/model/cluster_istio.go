package model

import (
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	uuid "github.com/satori/go.uuid"
)

type ClusterIstio struct {
	common.BaseModel
	ID        string `json:"-" gorm:"type:varchar(64)"`
	Name      string `json:"name"`
	ClusterID string `json:"cluster_id"`
	Version   string `json:"version"`
	Describe  string `json:"describe"`
	Status    string `json:"status"`
	Message   string `json:"message" gorm:"type:text(65535)"`
	Vars      string `json:"-" gorm:"type:text(65535)"`
}

func (c *ClusterIstio) BeforeCreate() (err error) {
	c.ID = uuid.NewV4().String()
	return nil
}
