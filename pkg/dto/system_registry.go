package dto

import "github.com/KubeOperator/KubeOperator/pkg/model"

type SystemRegistry struct {
	model.SystemRegistry
}

type SystemRegistryCreate struct {
	model.SystemRegistry
	Hostname     string `json:"hostname" validate:"koip,required"`
	Protocol     string `json:"protocol" validate:"oneof=http https"`
	Architecture string `json:"architecture" validate:"oneof=x86_64 aarch64"`
}

type SystemRegistryUpdate struct {
	ID           string `json:"id" validate:"required"`
	Hostname     string `json:"hostname" validate:"koip,required"`
	Protocol     string `json:"protocol" validate:"oneof=http https"`
	Architecture string `json:"architecture" validate:"oneof=x86_64 aarch64"`
}

type SystemRegistryBatchOp struct {
	Operation string           `json:"operation" validate:"required"`
	Items     []SystemRegistry `json:"items" validate:"required"`
}
