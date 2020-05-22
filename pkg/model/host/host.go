package host

import (
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	"github.com/KubeOperator/KubeOperator/pkg/model/credential"
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

func (h Host) TableName() string {
	return "ko_host"
}
