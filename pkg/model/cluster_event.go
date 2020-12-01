package model

import (
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	uuid "github.com/satori/go.uuid"
)

type ClusterEvent struct {
	common.BaseModel
	ID        string `json:"-"`
	UID       string `json:"uid"`
	Message   string `json:"message"`
	Kind      string `json:"kind"`
	Component string `json:"component"`
	Host      string `json:"host"`
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	Reason    string `json:"reason"`
	Detail    string `json:"detail"`
	ClusterID string `json:"clusterId"`
}

func (c *ClusterEvent) BeforeCreate() error {
	c.ID = uuid.NewV4().String()
	return nil
}
