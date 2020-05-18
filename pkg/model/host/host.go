package host

import (
	"ko3-gin/pkg/model/common"
	"ko3-gin/pkg/model/credential"
)

type Host struct {
	common.BaseModel
	credential.Credential
	CredentialId string
	Volumes      []Volume
}

type Volume struct {
	common.BaseModel
	size string
}
