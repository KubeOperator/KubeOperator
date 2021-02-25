package model

import (
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	uuid "github.com/satori/go.uuid"
)

type SystemRegistry struct {
	common.BaseModel
	ID               string `json:"id" gorm:"type:varchar(64)"`
	RegistryHostname string `json:"registry_hostname" gorm:"type:varchar(256);not null;unique"`
	RegistryProtocol string `json:"registry_protocol" gorm:"type:varchar(256);not null;"`
	Architecture     string `json:"architecture" gorm:"type:varchar(256);not null;"`
}

func (s *SystemRegistry) BeforeCreate() (err error) {
	s.ID = uuid.NewV4().String()
	return err
}
