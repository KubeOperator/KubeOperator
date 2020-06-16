package dto

import "github.com/KubeOperator/KubeOperator/pkg/model"

type Host struct {
	model.Host
}

type HostCreate struct {
	Name         string `json:"name" binding:"required"`
	Ip           string `json:"ip" binding:"required"`
	Port         int    `json:"port" binding:"required"`
	CredentialID string `json:"credentialId" binding:"required"`
}
