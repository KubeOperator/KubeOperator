package host

import (
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	"github.com/KubeOperator/KubeOperator/pkg/model/credential"
	uuid "github.com/satori/go.uuid"
)

const (
	Running      string = "Running"
	Waring       string = "Warning"
	Disconnect   string = "DisConnect"
	SshError     string = "SshError"
	AnsibleError string = "AnsibleError"
)

type Host struct {
	common.BaseModel
	Credential   credential.Credential
	ID           string
	Name         string `gorm:"not null;unique"`
	Memory       int
	CpuCore      int
	Os           string
	OsVersion    string
	GpuNum       int
	GpuInfo      string
	Ip           string
	Port         int
	CredentialID string
	Status       string
	NodeID       string
	Volumes      []Volume
}

type Volume struct {
	common.BaseModel
	size string
	name string
}

func (h *Host) BeforeCreate() (err error) {
	h.ID = uuid.NewV4().String()
	return err
}

func (h Host) TableName() string {
	return "ko_host"
}
