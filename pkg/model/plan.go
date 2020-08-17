package model

import (
	"errors"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	uuid "github.com/satori/go.uuid"
)

var (
	DeleteFailedByProject = "DELETE_FAILED_BY_PROJECT"
)

type Plan struct {
	common.BaseModel
	ID             string `json:"id" gorm:"type:varchar(64)"`
	Name           string `json:"name" gorm:"type:varchar(64)"`
	RegionID       string `json:"regionId" grom:"type:varchar(64)"`
	DeployTemplate string `json:"deployTemplate" grom:"type:varchar(64)"`
	Vars           string `json:"vars" gorm:"type text(65535)"`
	Zones          []Zone `json:"-" gorm:"many2many:ko_plan_zones"`
	Region         Region `json:"-"`
}

func (p *Plan) BeforeCreate() (err error) {
	p.ID = uuid.NewV4().String()
	return err
}

func (p Plan) TableName() string {
	return "ko_plan"
}

func (p *Plan) BeforeDelete() (err error) {
	var PlanResources []ProjectResource
	err = db.DB.Where(ProjectResource{ResourceId: p.ID}).Find(&PlanResources).Error
	if err != nil {
		return err
	}
	if len(PlanResources) > 0 {
		return errors.New(DeleteFailedByProject)
	}
	return err
}
