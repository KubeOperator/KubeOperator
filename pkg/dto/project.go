package dto

import "github.com/KubeOperator/KubeOperator/pkg/model"

type Project struct {
	model.Project
	Clusters interface{} `json:"clusters"`
}

type ProjectCreate struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}

type ProjectUpdate struct {
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
