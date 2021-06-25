package model

import (
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	uuid "github.com/satori/go.uuid"
)

type SystemRegistry struct {
	common.BaseModel
	ID                 string `json:"id" gorm:"type:varchar(64)"`
	Hostname           string `json:"hostname" gorm:"type:varchar(256);not null;unique"`
	Protocol           string `json:"protocol" gorm:"type:varchar(256);not null;"`
	Architecture       string `json:"architecture" gorm:"type:varchar(256);not null;"`
	RepoPort           int    `json:"repoPort" gorm:"type:int(64)"`
	RegistryPort       int    `json:"registryPort" gorm:"type:int(64)"`
	RegistryHostedPort int    `json:"registryHostedPort" gorm:"type:int(64)"`
}

func (s *SystemRegistry) BeforeCreate() (err error) {
	s.ID = uuid.NewV4().String()
	return err
}
