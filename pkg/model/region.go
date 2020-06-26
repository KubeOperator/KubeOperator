package model

import "github.com/KubeOperator/KubeOperator/pkg/model/common"

type Region struct {
	common.BaseModel
	ID         string `json:"id" gorm:"type:varchar(64)"`
	Name       string `json:"name" gorm:"type:varchar(256);not null;unique"`
	Datacenter string `json:"datacenter" gorm:"type:varchar(64)"`
}
