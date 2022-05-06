package model

import (
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	uuid "github.com/satori/go.uuid"
)

type KubepiBind struct {
	common.BaseModel
	ID           string `json:"id" gorm:"type:varchar(64)"`
	SourceType   string `json:"sourceType" gorm:"type:varchar(64)"`
	Source       string `json:"source" gorm:"type:varchar(64)"`
	BindUser     string `json:"bindUser" gorm:"type:varchar(64)"`
	BindPassword string `json:"bindPassword" gorm:"type:varchar(64)"`
}

func (k *KubepiBind) BeforeCreate() (err error) {
	k.ID = uuid.NewV4().String()
	return err
}
