package dto

import "github.com/KubeOperator/KubeOperator/pkg/model"

type SystemRegistry struct {
	model.SystemRegistry
}

type SystemRegistryCreate struct {
	model.SystemRegistry
	RegistryHostname string `json:"registry_hostname" validate:"required"`
	RegistryProtocol string `json:"registry_protocol" validate:"required"`
	Architecture     string `json:"architecture" validate:"required"`
}

type SystemRegistryUpdate struct {
	registry map[string]string `json:"vars" validate:"required"`
}

type SystemRegistryResult struct {
	registry map[string]string `json:"vars" validate:"required"`
}
