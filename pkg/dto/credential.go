package dto

import "github.com/KubeOperator/KubeOperator/pkg/model"

type Credential struct {
	model.Credential
}

type CredentialPage struct {
	Items []Credential `json:"items"`
	Total int          `json:"total"`
}

type CredentialCreate struct {
	Name       string `json:"name" validate:"required"`
	Username   string `json:"username" validate:"required"`
	Password   string
	PrivateKey string
	Type       string `json:"type" validate:"required"`
}

type CredentialUpdate struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	PrivateKey string `json:"privateKey"`
	Type       string `json:"type" validate:"required"`
}

type CredentialBatchOp struct {
	Operation string       `json:"operation" validate:"required"`
	Items     []Credential `json:"items" validate:"required"`
}
