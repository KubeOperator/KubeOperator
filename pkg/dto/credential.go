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
	Name       string `json:"name" validate:"koname,required,max=30"`
	Username   string `json:"username" validate:"required"`
	Password   string `json:"password" validate:"-"`
	PrivateKey string `json:"privateKey" validate:"-"`
	Type       string `json:"type" validate:"required"`
}

type CredentialUpdate struct {
	ID         string `json:"id" validate:"required"`
	Name       string `json:"name" validate:"koname,required,max=30"`
	Username   string `json:"username" validate:"required"`
	Password   string `json:"password" validate:"-"`
	PrivateKey string `json:"privateKey" validate:"-"`
	Type       string `json:"type" validate:"required"`
}

type CredentialBatchOp struct {
	Operation string       `json:"operation" validate:"required"`
	Items     []Credential `json:"items" validate:"required"`
}
