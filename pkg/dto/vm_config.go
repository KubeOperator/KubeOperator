package dto

import "github.com/KubeOperator/KubeOperator/pkg/model"

type VmConfig struct {
	model.VmConfig
}

type VmConfigOp struct {
	Operation string     `json:"operation"`
	Items     []VmConfig `json:"items"`
}

type VmConfigCreate struct {
	Name     string `json:"name" validate:"required"`
	Provider string `json:"provider"`
	Cpu      int    `json:"cpu" validate:"required"`
	Memory   int    `json:"memory" validate:"required"`
}

type VmConfigUpdate struct {
	Name     string `json:"name" validate:"required"`
	Provider string `json:"provider"`
	Cpu      int    `json:"cpu" validate:"required"`
	Memory   int    `json:"memory" validate:"required"`
}
