package dto

import "github.com/KubeOperator/KubeOperator/pkg/model"

type SystemRegistry struct {
	model.SystemRegistry
}

type SystemRegistryCreate struct {
	model.SystemRegistry
	Hostname     string `json:"hostname" validate:"required"`
	Protocol     string `json:"protocol" validate:"required"`
	Architecture string `json:"architecture" validate:"required"`
}

type SystemRegistryUpdate struct {
	Registry map[string]string `json:"vars" validate:"required"`
}

type SystemRegistryResult struct {
	Registry map[string]string `json:"vars" validate:"required"`
}
