package dto

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/model"
)

type Plan struct {
	model.Plan
	Zones []Zone `json:"zones"`
}

type PlanCreate struct {
	Name           string      `json:"name" validate:"required"`
	RegionId       string      `json:"regionId" validate:"required"`
	Zones          []string    `json:"zones" validate:"required"`
	PlanVars       interface{} `json:"planVars" validate:"required"`
	DeployTemplate string      `json:"deployTemplate" validate:"required"`
}

type PlanOp struct {
	Operation string `json:"operation" validate:"required"`
	Items     []Plan `json:"items" validate:"required"`
}

type PlanVmConfig struct {
	Name   string            `json:"name"`
	Config constant.VmConfig `json:"config"`
}
