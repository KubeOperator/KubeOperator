package dto

import "github.com/KubeOperator/KubeOperator/pkg/model"

type User struct {
	model.User
}

type UserCreate struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	IsAdmin  bool   `json:"isAdmin" binding:"required"`
}

type UserPage struct {
	Items []User `json:"items"`
	Total int    `json:"total"`
}

type UserUpdate struct {
	ID       string `json:"id"`
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	IsActive bool   `json:"isActive"`
	IsAdmin  bool   `json:"isAdmin" binding:"required"`
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
