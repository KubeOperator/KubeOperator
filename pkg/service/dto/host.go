package dto

import "github.com/KubeOperator/KubeOperator/pkg/model"

type Host struct {
	model.Host
}

type HostCreate struct {
	Name         string `json:"name" validate:"required"`
	Ip           string `json:"ip" validate:"required"`
	Port         int    `json:"port" validate:"required"`
	CredentialID string `json:"credentialId" validate:"required"`
}

type HostPage struct {
	Items []Host `json:"items"`
	Total int    `json:"total"`
}

type HostOp struct {
	Operation string `json:"operation" validate:"required"`
	Items     []Host `json:"items" validate:"required"`
}
