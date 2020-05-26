package host

import (
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	"github.com/KubeOperator/KubeOperator/pkg/model/credential"
	uuid "github.com/satori/go.uuid"
)

type Host struct {
	common.BaseModel
	Credential   credential.Credential
	ID           string
	Name         string
	Memory       string
	CpuCore      int
	Os           string
	OsVersion    string
	GpuNum       int
	GpuInfo      string
	Ip           string
	Port         int
	CredentialId string
	ClusterId    string
	nodeId       string
	Status       string
	Volumes      []Volume
}

type Volume struct {
	common.BaseModel
	size string
}

func (h *Host) BeforeCreate() (err error) {
	h.ID = uuid.NewV4().String()
	return err
}

func (h Host) TableName() string {
	return "ko_host"
}
