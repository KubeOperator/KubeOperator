package dto

import "github.com/KubeOperator/KubeOperator/pkg/model"

type Plan struct {
	model.Plan
}

type PlanCreate struct {
	Name string `json:"name" validate:"required"`
}

type PlanOp struct {
	Operation string `json:"operation" validate:"required"`
	Items     []Plan `json:"items" validate:"required"`
}
