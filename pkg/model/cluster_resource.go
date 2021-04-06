package model

import (
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	"github.com/google/uuid"
)

type ClusterResource struct {
	common.BaseModel
	ID           string  `json:"id"`
	ResourceType string  `json:"resourceType"`
	ResourceID   string  `json:"resourceId"`
	ClusterID    string  `json:"clusterId"`
	Cluster      Cluster `json:"-"`
}

func (c *ClusterResource) BeforeCreate() (err error) {
	c.ID = uuid.New().String()
	return err
}
