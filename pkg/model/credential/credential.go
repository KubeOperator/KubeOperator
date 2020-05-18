package credential

import "ko3-gin/pkg/model/common"

type Credential struct {
	common.BaseModel
	Username   string
	Password   string
	PrivateKey string
	Type       string
}
