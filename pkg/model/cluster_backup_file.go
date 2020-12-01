package model

import (
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	uuid "github.com/satori/go.uuid"
)

type ClusterBackupFile struct {
	common.BaseModel
	ID                      string                `json:"id"`
	Name                    string                `json:"name"`
	ClusterID               string                `json:"clusterId"`
	ClusterBackupStrategyID string                `json:"clusterBackupStrategyId"`
	Folder                  string                `json:"folder"`
	ClusterBackupStrategy   ClusterBackupStrategy `json:"-"`
	CLuster                 Cluster               `json:"-"`
}

func (c *ClusterBackupFile) BeforeCreate() error {
	c.ID = uuid.NewV4().String()
	return nil
}
