package model

import (
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

func (m *MultiClusterSyncClusterLog) BeforeCreate() error {
	m.ID = uuid.NewV4().String()
	return nil
}
