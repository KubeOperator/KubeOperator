package model

import (
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	uuid "github.com/satori/go.uuid"
)

type ClusterMultiClusterRepository struct {
	common.BaseModel
	ID                       string `json:"-"`
	ClusterID                string `json:"clusterId"`
	MultiClusterRepositoryID string `json:"multiClusterRepositoryId"`
	Status                   string `json:"status"`
	Message                  string `json:"message"`
}

func (c *ClusterMultiClusterRepository) BeforeCreate() (err error) {
	c.ID = uuid.NewV4().String()
	return nil
}
