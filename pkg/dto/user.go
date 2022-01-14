package dto

import "github.com/KubeOperator/KubeOperator/pkg/model"

type User struct {
	model.User
}

type UserCreate struct {
	Name     string `json:"name" validate:"required,max=30"`
	Email    string `json:"email" validate:"-"`
	Password string `json:"password" validate:"kopassword,required"`
	IsAdmin  bool   `json:"isAdmin" validate:"-"`
}

type UserPage struct {
	Items []User `json:"items"`
	Total int    `json:"total"`
}

type UserUpdate struct {
	ID       string `json:"id"`
	Name     string `json:"name" validate:"required,max=30"`
	Email    string `json:"email" validate:"-"`
	IsActive bool   `json:"isActive" validate:"-"`
	IsAdmin  bool   `json:"isAdmin" validate:"-"`
}

type UserOp struct {
	Operation string `json:"operation"`
	Items     []User `json:"items"`
}

type UserChangePassword struct {
	ID       string `json:"id" validate:"-"`
	Name     string `json:"name" validate:"required,min=6"`
	Password string `json:"password" validate:"kopassword,required,max=30,min=6"`
	Original string `json:"original" validate:"kopassword,required,max=30,min=6"`
}

type UserForgotPassword struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
}
