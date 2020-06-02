package cluster

import (
	commonModel "github.com/KubeOperator/KubeOperator/pkg/model/common"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

type Status struct {
	commonModel.BaseModel
	ID         string
	Version    string
	Message    string `gorm:"type:text(65535)"`
	Phase      string
	Conditions []Condition `gorm:"save_associations:false"`
}

func (s *Status) BeforeCreate() (err error) {
	s.ID = uuid.NewV4().String()
	return nil
}

func (s Status) BeforeDelete(scope *gorm.Scope) error {
	return scope.DB().
		Where(Condition{StatusID: s.ID}).
		Delete(Condition{}).Error
}

func (s Status) TableName() string {
	return "ko_cluster_status"
}
