package model

import (
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	uuid "github.com/satori/go.uuid"
)

type F5Setting struct {
	common.BaseModel
	ID        string `json:"id" gorm:"type:varchar(64)"`
	ClusterID string `json:"clusterID" gorm:"type:varchar(64)"`
	URL       string `json:"url" gorm:"type:varchar(64)"`
	User      string `json:"user" gorm:"type:varchar(64)"`
	Partition string `json:"partition" gorm:"type:varchar(64)" `
	PublicIP  string `json:"publicIP" gorm:"type:varchar(64)"`
	Status    string `json:"status" gorm:"type:varchar(64)"`
}

func (s *F5Setting) BeforeCreate() (err error) {
	s.ID = uuid.NewV4().String()
	return err
}
