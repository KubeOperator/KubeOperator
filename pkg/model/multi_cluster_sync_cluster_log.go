package model

import (
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	uuid "github.com/satori/go.uuid"
)

type MultiClusterSyncClusterLog struct {
	common.BaseModel
	ID                    string `json:"-"`
	Status                string `json:"status"`
	Message               string `json:"message"`
	MultiClusterSyncLogID string `json:"multiClusterSyncLogId"`
	ClusterID             string `json:"clusterId"`
}

func (m *MultiClusterSyncClusterLog) BeforeDelete() error {
	var mls []MultiClusterSyncClusterResourceLog
	if err := db.DB.Where(MultiClusterSyncClusterResourceLog{
		MultiClusterSyncClusterLogID: m.ID,
	}).Find(&mls).Error; err != nil {
		return err
	}
	tx := db.DB.Begin()
	for _,m := range mls {
		if err:=db.DB.Delete(&m).Error;err!=nil{
			tx.Rollback();
			return err
		}
	}
	tx.Commit()
	return nil
}

func (m *MultiClusterSyncClusterLog) BeforeCreate() error {
	m.ID = uuid.NewV4().String()
	return nil
}
