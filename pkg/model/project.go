package model

import (
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	uuid "github.com/satori/go.uuid"
)

type Project struct {
	common.BaseModel
	ID          string    `json:"id" gorm:"type:varchar(64)"`
	Name        string    `json:"name" gorm:"type:varchar(64);not null;unique"`
	Description string    `json:"description" gorm:"type:varchar(128)"`
	Clusters    []Cluster `json:"clusters"`
}

func (p *Project) BeforeCreate() (err error) {
	p.ID = uuid.NewV4().String()
	return err
}

func (p Project) TableName() string {
	return "ko_project"
}
