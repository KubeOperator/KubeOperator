package dto

import "github.com/KubeOperator/KubeOperator/pkg/model"

type ProjectMember struct {
	model.ProjectMember
	UserName string `json:"userName"`
}

type ProjectMemberOP struct {
	Operation string          `json:"operation" validate:"required"`
	Items     []ProjectMember `json:"items" validate:"required"`
}
