package credential

import "github.com/KubeOperator/KubeOperator/pkg/model/common"

type Credential struct {
	common.BaseModel
	Username   string
	Password   string
	PrivateKey string
	Type       string
}
