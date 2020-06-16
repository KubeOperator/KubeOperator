package dto

import "github.com/KubeOperator/KubeOperator/pkg/model"

type Credential struct {
	model.Credential
}

type CredentialCreate struct {
	Name       string `json:"name" binding:"required"`
	Username   string `json:"username" binding:"required"`
	Password   string
	PrivateKey string
	Type       string `json:"type" binding:"required"`
}
