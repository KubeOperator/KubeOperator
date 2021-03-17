package dto

import "github.com/KubeOperator/KubeOperator/pkg/model"

type User struct {
	model.User
	Role   string `json:"role"`
	Status string `json:"status"`
}

type UserCreate struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role" binding:"required"`
}

type UserPage struct {
	Items []User `json:"items"`
	Total int    `json:"total"`
}

type UserUpdate struct {
	Email          string `json:"email"`
	Password       string `json:"password"`
	Role           string `json:"role"`
	Language       string `json:"language"`
	CurrentProject string `json:"currentProject"`
}

type UserOp struct {
	Operation string `json:"operation"`
	Items     []User `json:"items"`
}

type UserChangePassword struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Original string `json:"original"`
}

type UserForgotPassword struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
}
