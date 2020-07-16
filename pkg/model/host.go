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
	ID           string     `json:"id" gorm:"type:varchar(64)"`
	Name         string     `json:"name" gorm:"type:varchar(256);not null;unique"`
	Memory       int        `json:"memory" gorm:"type:int(64)"`
	CpuCore      int        `json:"cpuCore" gorm:"type:int(64)"`
	Os           string     `json:"os" gorm:"type:varchar(64)"`
	OsVersion    string     `json:"osVersion" gorm:"type:varchar(64)"`
	GpuNum       int        `json:"gpuNum" gorm:"type:int(64)"`
	GpuInfo      string     `json:"gpuInfo" gorm:"type:varchar(128)"`
	Ip           string     `json:"ip" gorm:"type:varchar(128);not null;unique"`
	Port         int        `json:"port" gorm:"type:varchar(64)"`
	CredentialID string     `json:"credentialId" gorm:"type:varchar(64)"`
	Status       string     `json:"status" gorm:"type:varchar(64)"`
	ClusterID    string     `json:"clusterId" gorm:"type:varchar(64)"`
	ZoneID       string     `json:"zoneId"`
	Zone         Zone       `gorm:"save_associations:false" json:"_"`
	Volumes      []Volume   `gorm:"save_associations:false" json:"volumes"`
	Credential   Credential `gorm:"save_associations:false" json:"credential"`
	Cluster      Cluster    `gorm:"save_associations:false" json:"cluster"`
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
