package dto

import "github.com/KubeOperator/KubeOperator/pkg/model"

type ProjectResource struct {
	model.ProjectResource
	ResourceName string `json:"resourceName"`
}

type ProjectResourceOp struct {
	Operation string            `json:"operation" validate:"required"`
	Items     []ProjectResource `json:"items" validate:"required"`
}
