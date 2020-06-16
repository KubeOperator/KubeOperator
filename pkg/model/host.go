package model

import (
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	uuid "github.com/satori/go.uuid"
)

const (
	Running      string = "Running"
	Warning      string = "Warning"
	Disconnect   string = "DisConnect"
	SshError     string = "SshError"
	AnsibleError string = "AnsibleError"
)

type Host struct {
	common.BaseModel
	Credential   Credential
	ID           string
	Name         string `gorm:"not null;unique"`
	Memory       int
	CpuCore      int
	Os           string
	OsVersion    string
	GpuNum       int
	GpuInfo      string
	Ip           string `gorm:"not null;unique"`
	Port         int
	CredentialID string
	Status       string
	NodeID       string
	Volumes      []Volume
}

func (h *Host) BeforeCreate() (err error) {
	h.ID = uuid.NewV4().String()
	return err
}

func (h Host) TableName() string {
	return "ko_host"
}
