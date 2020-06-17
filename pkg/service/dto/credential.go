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
	Name       string `json:"name" binding:"required"`
	Username   string `json:"username" binding:"required"`
	Password   string
	PrivateKey string
	Type       string `json:"type" binding:"required"`
}

type CredentialUpdate struct {
	ID         string `json:"id" binding:"required"`
	Name       string `json:"name" binding:"required"`
	Username   string `json:"username" binding:"required"`
	Password   string
	PrivateKey string
	Type       string `json:"type" binding:"required"`
}

type CredentialBatchOp struct {
	Operation string       `json:"operation"`
	Items     []Credential `json:"items"`
}
