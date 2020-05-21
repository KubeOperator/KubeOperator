package credential

import (
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
)

type Credential struct {
	common.BaseModel
	ID         string
	Name       string
	Username   string
	Password   string
	PrivateKey string
	Type       string
}

func (c Credential) TableName() string {
	return "ko_credential"
}
