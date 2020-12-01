package model

import (
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	uuid "github.com/satori/go.uuid"
)

type Message struct {
	common.BaseModel
	ID        string  `json:"-"`
	Title     string  `json:"title"`
	Content   string  `json:"content"`
	Type      string  `json:"type"`
	Level     string  `json:"level"`
	ClusterID string  `json:"clusterId"`
	Cluster   Cluster `json:"-"`
}

func (m *Message) BeforeCreate() error {
	m.ID = uuid.NewV4().String()
	return nil
}
