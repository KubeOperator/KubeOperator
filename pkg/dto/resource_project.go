package dto

import "github.com/KubeOperator/KubeOperator/pkg/model"

type ProjectResource struct {
	model.ProjectResource
}

type ProjectResourceCreate struct {
	ProjectID    string   `json:"projectId"`
	ResourceType string   `json:"resourceType"`
	ResourceIds  []string `json:"resourceIds"`
}

type ProjectResourceOp struct {
	Operation string            `json:"operation" validate:"required"`
	Items     []ProjectResource `json:"items" validate:"required"`
}
