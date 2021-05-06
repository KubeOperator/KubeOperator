package model

import (
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	uuid "github.com/satori/go.uuid"
)

type ProjectResource struct {
	common.BaseModel
	ID           string  `json:"id" gorm:"type:varchar(64)"`
	ResourceType string  `json:"resourceType" gorm:"type:varchar(128)" validate:"required"`
	ResourceID   string  `json:"resourceId" gorm:"type:varchar(64)" validate:"required"`
	ProjectID    string  `json:"projectId" gorm:"type:varchar(64)" validate:"required"`
	Project      Project `json:"-"`
}

func (p *ProjectResource) BeforeCreate() (err error) {
	p.ID = uuid.NewV4().String()
	return err
}
