package model

import (
	"errors"

	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	"github.com/jinzhu/gorm"
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
	Zones          []Zone `json:"-" gorm:"many2many:plan_zones"`
	Region         Region `json:"-"`
}

func (p *Plan) BeforeCreate() (err error) {
	p.ID = uuid.NewV4().String()
	return err
}

func (p *Plan) BeforeDelete(tx *gorm.DB) (err error) {
	var PlanResources []ProjectResource
	if err := tx.Where(ProjectResource{ResourceID: p.ID}).Find(&PlanResources).Error; err != nil {
		return err
	}
	if len(PlanResources) > 0 {
		return errors.New(DeleteFailedByProject)
	}
	return nil
}
