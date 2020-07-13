package model

import (
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	uuid "github.com/satori/go.uuid"
)

type Plan struct {
	common.BaseModel
	ID             string `json:"id" gorm:"type:varchar(64)"`
	Name           string `json:"name" gorm:"type:varchar(64)"`
	RegionID       string `json:"regionId" grom:"type:varchar(64)"`
	DeployTemplate string `json:"deployTemplate" grom:"type:varchar(64)"`
	Vars           string `json:"vars" gorm:"type text(0)"`
	Zones          []Zone `json:"zones"`
	Region         Region `json:"region"`
}

func (p *Plan) BeforeCreate() (err error) {
	p.ID = uuid.NewV4().String()
	return err
}

func (p Plan) TableName() string {
	return "ko_plan"
}
