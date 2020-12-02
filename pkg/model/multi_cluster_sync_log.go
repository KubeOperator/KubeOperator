package model

import (
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	uuid "github.com/satori/go.uuid"
)

type MultiClusterSyncLog struct {
	common.BaseModel
	ID                       string `json:"-"`
	Status                   string `json:"status"`
	Message                  string `json:"message"`
	MultiClusterRepositoryID string `json:"multiClusterRepositoryId"`
}

func (m *MultiClusterSyncLog) BeforeCreate() error {
	m.ID = uuid.NewV4().String()
	return nil
}
