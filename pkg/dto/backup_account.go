package dto

import "github.com/KubeOperator/KubeOperator/pkg/model"

type BackupAccount struct {
	model.BackupAccount
	CredentialVars interface{} `json:"credentialVars"`
}

type BackupAccountCreate struct {
	Name           string      `json:"name" validate:"required"`
	CredentialVars interface{} `json:"credentialVars" validate:"required"`
	Region         string      `json:"region" validate:"required"`
	Type           string      `json:"type" validate:"required"`
}

type BackupAccountOp struct {
	Operation string          `json:"operation" validate:"required"`
	Items     []BackupAccount `json:"items" validate:"required"`
}

type BackupAccountUpdate struct {
	Name           string      `json:"name" validate:"required"`
	CredentialVars interface{} `json:"credentialVars" validate:"required"`
	Region         string      `json:"region" validate:"required"`
	Type           string      `json:"type" validate:"required"`
}
