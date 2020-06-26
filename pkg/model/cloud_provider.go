package model

import (
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	uuid "github.com/satori/go.uuid"
)

type CloudProvider struct {
	common.BaseModel
	ID   string `json:"id" gorm:"type:varchar(64)"`
	Name string `json:"name" gorm:"type:varchar(64)"`
	Vars string `json:"vars" gorm:"type longtext(0)"`
}

func (c *CloudProvider) BeforeCreate() (err error) {
	c.ID = uuid.NewV4().String()
	return err
}

func (c CloudProvider) TableName() string {
	return "ko_cloud_provider"
}
