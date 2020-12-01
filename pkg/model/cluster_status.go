package model

import (
	"github.com/KubeOperator/KubeOperator/pkg/db"
	commonModel "github.com/KubeOperator/KubeOperator/pkg/model/common"
	uuid "github.com/satori/go.uuid"
)

type ClusterStatus struct {
	commonModel.BaseModel
	ID                      string
	Message                 string                   `json:"message" gorm:"type:text(65535)"`
	Phase                   string                   `json:"phase"`
	PrePhase                string                   `json:"prePhase"`
	ClusterStatusConditions []ClusterStatusCondition `json:"conditions" gorm:"save_associations:false"`
}

func (s *ClusterStatus) BeforeCreate() (err error) {
	s.ID = uuid.NewV4().String()
	return nil
}

func (s ClusterStatus) BeforeDelete() (err error) {
	status := ClusterStatus{ID: s.ID}
	if err := db.DB.
		First(&status).
		Related(&status.ClusterStatusConditions).Error; err != nil {
		return err
	}
	tx := db.DB.Begin()
	if len(status.ClusterStatusConditions) > 0 {
		if err := db.DB.
			Where(ClusterStatusCondition{ClusterStatusID: status.ID}).
			Delete(ClusterStatusCondition{}).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	return nil
}
