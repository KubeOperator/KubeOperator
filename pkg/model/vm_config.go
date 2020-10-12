package model

import (
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	uuid "github.com/satori/go.uuid"
)

type VmConfig struct {
	common.BaseModel
	ID       string `json:"-"`
	Name     string `json:"name"`
	Cpu      int    `json:"cpu"`
	Memory   int    `json:"memory"`
	Disk     int    `json:"disk"`
	Provider string `json:"provider"`
}

func (v *VmConfig) BeforeCreate() error {
	v.ID = uuid.NewV4().String()
	return nil
}

func (v VmConfig) TableName() string {
	return "ko_vm_config"
}
