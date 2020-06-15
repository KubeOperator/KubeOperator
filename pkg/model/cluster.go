package model

import (
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

type Cluster struct {
	common.BaseModel
	ID       string
	Name     string
	SpecID   string
	SecretID string
	StatusID string
	Spec     ClusterSpec   `gorm:"save_associations:false"`
	Secret   ClusterSecret `gorm:"save_associations:false"`
	Status   ClusterStatus `gorm:"save_associations:false"`
}

func (c Cluster) TableName() string {
	return "ko_cluster"
}

func (c *Cluster) BeforeCreate(scope *gorm.Scope) error {
	c.ID = uuid.NewV4().String()
	return nil
}
