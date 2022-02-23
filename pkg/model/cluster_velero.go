package model

import (
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	uuid "github.com/satori/go.uuid"
)

type ClusterVelero struct {
	common.BaseModel
	ID                string `json:"id"`
	Cluster           string `json:"cluster"`
	BackupAccountName string `json:"backupAccountName"`
	Bucket            string `json:"bucket"`
	Endpoint          string `json:"endpoint"`
}

func (c *ClusterVelero) BeforeCreate() error {
	c.ID = uuid.NewV4().String()
	return nil
}
