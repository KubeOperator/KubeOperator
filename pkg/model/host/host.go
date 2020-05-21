package host

import (
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	"github.com/KubeOperator/KubeOperator/pkg/model/credential"
)

type Host struct {
	common.BaseModel
	credential.Credential
	ID           string
	Name         string
	Ip           string
	User         string
	Password     string
	Port         int
	CredentialId string
	Status       string
	Volumes      []Volume
}

type Volume struct {
	common.BaseModel
	size string
}

func (h Host) TableName() string {
	return "ko_credential"
}
