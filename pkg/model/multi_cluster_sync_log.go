package model

import (
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	uuid "github.com/satori/go.uuid"
)

type MultiClusterSyncLog struct {
	common.BaseModel
	ID                       string `json:"id"`
	Status                   string `json:"status"`
	Message                  string `json:"message"`
	GitCommitId              string `json:"gitCommitId"`
	MultiClusterRepositoryID string `json:"multiClusterRepositoryId"`
}

func (m *MultiClusterSyncLog) BeforeDelete() error {
	var mls []MultiClusterSyncClusterLog
	if err := db.DB.Where(MultiClusterSyncClusterLog{
		MultiClusterSyncLogID: m.ID,
	}).Find(&mls).Error; err != nil {
		return err
	}
	tx := db.DB.Begin()
	for _,m := range mls {
		if err := db.DB.Delete(&m).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	return nil
}

func (m *MultiClusterSyncLog) BeforeCreate() error {
	m.ID = uuid.NewV4().String()
	return nil
}
