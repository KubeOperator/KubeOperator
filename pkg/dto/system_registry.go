package dto

import "github.com/KubeOperator/KubeOperator/pkg/model"

type SystemRegistry struct {
	model.SystemRegistry
}

type SystemRegistryCreate struct {
	registry map[string]string `json:"vars" validate:"required"`
}

type SystemRegistryUpdate struct {
	registry map[string]string `json:"vars" validate:"required"`
}

type SystemRegistryResult struct {
	registry map[string]string `json:"vars" validate:"required"`
}
