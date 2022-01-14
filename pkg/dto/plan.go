package dto

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/model"
)

type Plan struct {
	model.Plan
	PlanVars   interface{} `json:"planVars"`
	RegionName string      `json:"regionName"`
	ZoneNames  []string    `json:"zoneNames"`
}

type PlanCreate struct {
	Name           string      `json:"name" validate:"koname,required"`
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
