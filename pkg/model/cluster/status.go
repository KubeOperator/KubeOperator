package cluster

import (
	"github.com/KubeOperator/KubeOperator/pkg/db"
	commonModel "github.com/KubeOperator/KubeOperator/pkg/model/common"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

type Status struct {
	commonModel.BaseModel
	ID           string
	Message      string `gorm:"type:text(65535)"`
	Phase        string
	Conditions   []Condition `gorm:"save_associations:false"`
}

func (s *Status) BeforeCreate() (err error) {
	s.ID = uuid.NewV4().String()
	return nil
}

func (s *Status) AfterSave() error {
	if len(s.Conditions) > 0 {
		for i, _ := range s.Conditions {
			s.Conditions[i].StatusID = s.ID
			if db.DB.NewRecord(s.Conditions[i]) {
				s.Conditions[i].ID = uuid.NewV4().String()
				if err := db.DB.
					Create(&(s.Conditions[i])).Error; err != nil {
					return err
				}
			} else {
				if err := db.DB.
					Save(&(s.Conditions[i])).Error; err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (s *Status) ClearConditions() error {
	if len(s.Conditions) > 0 {
		err := db.DB.Where(Condition{StatusID: s.ID}).
			Delete(Condition{}).Error
		if err != nil {
			return err
		}
		s.Conditions = []Condition{}
	}
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
