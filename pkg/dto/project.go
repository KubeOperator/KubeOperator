package dto

import "github.com/KubeOperator/KubeOperator/pkg/model"

type Project struct {
	model.Project
}

type ProjectCreate struct {
	Name        string `json:"name" validate:"koname,required,max=30"`
	Description string `json:"description" validate:"max=255"`
	UserName    string `json:"userName"`
}

type ProjectUpdate struct {
	ID          string `json:"id"`
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}
type ProjectPage struct {
	Items []Project `json:"items"`
	Total int       `json:"total"`
}

type ProjectOp struct {
	Operation string    `json:"operation" validate:"required"`
	Items     []Project `json:"items" validate:"required"`
}
