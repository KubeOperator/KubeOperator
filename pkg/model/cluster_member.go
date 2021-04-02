package model

import (
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	uuid "github.com/satori/go.uuid"
)

type ClusterMember struct {
	common.BaseModel
	ID        string `json:"-" gorm:"type:varchar(64)"`
	ClusterID string `json:"clusterId" gorm:"type:varchar(64)"`
	UserID    string `json:"userId" gorm:"type:varchar(64)"`
	Role      string `json:"role" gorm:"type:varchar(64)"`
	User      User   `json:"user" gorm:"save_associations:false"`
}

func (c *ClusterMember) BeforeCreate() (err error) {
	c.ID = uuid.NewV4().String()
	return err
}
