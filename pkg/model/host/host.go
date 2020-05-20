package host

import (
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	"github.com/KubeOperator/KubeOperator/pkg/model/credential"
)

type Host struct {
	common.BaseModel
	credential.Credential
	Ip       string
	User     string
	Password string
	Port     int
	CredentialId string
	Status       string
	Volumes      []Volume
}

type Volume struct {
	common.BaseModel
	size string
}
