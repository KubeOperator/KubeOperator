package model

import (
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	uuid "github.com/satori/go.uuid"
)

type StorageProvisionerDic struct {
	common.BaseModel
	ID           string `json:"id"`
	Name         string `json:"name"`
	Version      string `json:"version"`
	Architecture string `json:"architecture"`
	Vars         string `json:"-"    gorm:"type:text(65535)"`
}

func (c *StorageProvisionerDic) BeforeCreate() (err error) {
	c.ID = uuid.NewV4().String()
	return nil
}
