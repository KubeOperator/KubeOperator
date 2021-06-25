package dto

import "github.com/KubeOperator/KubeOperator/pkg/model"

type SystemRegistry struct {
	model.SystemRegistry
}

type SystemRegistryCreate struct {
	model.SystemRegistry
	Hostname           string `json:"hostname" validate:"required"`
	Protocol           string `json:"protocol" validate:"required"`
	Architecture       string `json:"architecture" validate:"required"`
	RepoPort           int    `json:"repoPort" validate:"required"`
	RegistryPort       int    `json:"registryPort" validate:"required"`
	RegistryHostedPort int    `json:"registryHostedPort" validate:"required"`
}

type SystemRegistryUpdate struct {
	ID       string `json:"id" validate:"required"`
	Hostname string `json:"hostname" validate:"required"`
	Protocol string `json:"protocol" validate:"required"`
}

type SystemRegistryDelete struct {
	Architecture string `json:"architecture" validate:"required"`
}

type SystemRegistryBatchOp struct {
	Operation string           `json:"operation" validate:"required"`
	Items     []SystemRegistry `json:"items" validate:"required"`
}
