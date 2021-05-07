package dto

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/model"
)

type Plan struct {
	model.Plan
	PlanVars interface{} `json:"planVars"`
	Region   string      `json:"region"`
	Zones    []string    `json:"zones"`
	Projects []string    `json:"projects"`
	Provider string      `json:"provider"`
}

type PlanCreate struct {
	Name           string      `json:"name" validate:"required"`
	Zones          []string    `json:"zones" validate:"required"`
	PlanVars       interface{} `json:"planVars" validate:"required"`
	DeployTemplate string      `json:"deployTemplate" validate:"required"`
	Projects       []string    `json:"projects" validate:"required"`
	Region         string      `json:"region" validate:"required"`
}

type PlanOp struct {
	Operation string `json:"operation" validate:"required"`
	Items     []Plan `json:"items" validate:"required"`
}

type PlanVmConfig struct {
	Name   string            `json:"name"`
	Config constant.VmConfig `json:"config"`
}

type PlanUpdate struct {
	PlanVars interface{} `json:"planVars" validate:"required"`
	Projects []string    `json:"projects" validate:"required"`
}
