package model

import (
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	uuid "github.com/satori/go.uuid"
	"time"
)

type MultiClusterRepository struct {
	common.BaseModel
	ID           string `json:"-"`
	Name         string
	Source       string
	Username     string
	Password     string
	Status       string
	Message      string
	Branch       string
	GitCommitId  string
	LastSyncTime time.Time
	SyncInterval int64
	GitTimeout   int64
	SyncEnable   bool
}

func (m *MultiClusterRepository) BeforeCreate() error {
	m.ID = uuid.NewV4().String()
	return nil
}

func (m *MultiClusterRepository) Pull() error {

}
