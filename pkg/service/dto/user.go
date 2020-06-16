package dto

import "github.com/KubeOperator/KubeOperator/pkg/model"

type User struct {
	model.User
}

type UserCreate struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}
