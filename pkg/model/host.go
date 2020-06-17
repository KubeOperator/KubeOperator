package model

import (
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	"github.com/KubeOperator/KubeOperator/pkg/util/encrypt"
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
	ID           string     `json:"id"`
	Name         string     `json:"name" gorm:"not null;unique"`
	Memory       int        `json:"memory"`
	CpuCore      int        `json:"cpuCore"`
	Os           string     `json:"os"`
	OsVersion    string     `json:"osVersion"`
	GpuNum       int        `json:"gpuNum"`
	GpuInfo      string     `json:"gpuInfo"`
	Ip           string     `json:"ip" gorm:"not null;unique"`
	Port         int        `json:"port"`
	CredentialID string     `json:"credentialId"`
	Status       string     `json:"status"`
	Volumes      []Volume   `json:"volumes"`
	Credential   Credential `json:"credential"`
}

func (h Host) GetHostPasswordAndPrivateKey() (string, []byte, error) {
	var err error = nil
	password := ""
	privateKey := []byte("")
	if "password" == h.Credential.Type {
		pwd, err := encrypt.StringDecrypt(h.Credential.Password)
		password = pwd
		if err != nil {
			return password, privateKey, err
		}
	}
	if "privateKey" == h.Credential.Type {
		privateKey = []byte(h.Credential.PrivateKey)
	}
	return password, privateKey, err
}

func (h *Host) BeforeCreate() (err error) {
	h.ID = uuid.NewV4().String()
	return err
}

func (h Host) TableName() string {
	return "ko_host"
}
