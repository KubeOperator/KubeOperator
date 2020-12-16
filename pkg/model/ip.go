package model

import (
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	uuid "github.com/satori/go.uuid"
)

type Ip struct {
	common.BaseModel
	ID       string `json:"id" gorm:"type:varchar(64)"`
	Ip       string `json:"ip" gorm:"type:varchar(255)"`
	Gateway  string `json:"gateway" gorm:"type:varchar(255)"`
	DNS1     string `json:"dns1" gorm:"type:varchar(255)"`
	DNS2     string `json:"dns2" gorm:"type:varchar(255)"`
	Status   string `json:"status" gorm:"type:varchar(255)"`
	IpPoolID string `json:"ipPoolId" gorm:"type:varchar(64)"`
}

func (i *Ip) BeforeCreate() (err error) {
	i.ID = uuid.NewV4().String()
	return err
}
