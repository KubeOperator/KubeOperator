package model

import (
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	uuid "github.com/satori/go.uuid"
)

type ClusterToolDetail struct {
	common.BaseModel
	ID           string `json:"-" gorm:"type:varchar(64)"`
	Name         string `json:"name"`
	Version      string `json:"version"`
	ChartVersion string `json:"chart_version"`
	Architecture string `json:"architecture"`
	Vars         string `json:"-"  gorm:"type:text(65535)"`
}

func (c *ClusterToolDetail) BeforeCreate() (err error) {
	c.ID = uuid.NewV4().String()
	return nil
}
