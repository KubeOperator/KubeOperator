package dto

import "github.com/KubeOperator/KubeOperator/pkg/model"

type ProjectMember struct {
	model.ProjectMember
	Email      string `json:"email"`
	UserStatus string `json:"userStatus"`
	UserName   string `json:"userName"`
}

type ProjectMemberOP struct {
	Operation string                `json:"operation" validate:"required"`
	Items     []ProjectMemberCreate `json:"items" validate:"required"`
}

type ProjectMemberCreate struct {
	ProjectName string `json:"projectName" validate:"required"`
	Role        string `json:"role" validate:"required"`
	Username    string `json:"username" validate:"required"`
}

type AddMemberResponse struct {
	Items []string `json:"items"`
}
