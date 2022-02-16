package dto

import "time"

type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	IsActive  bool      `json:"isActive"`
	Language  string    `json:"language"`
	IsAdmin   bool      `json:"isAdmin"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"createdAt"`
}

type UserCreate struct {
	Name     string `json:"name" validate:"koname,required,max=30"`
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
	IsActive bool   `json:"isActive" validate:"-"`
	IsAdmin  bool   `json:"isAdmin" validate:"-"`
}

type UserOp struct {
	Operation string `json:"operation"`
	Items     []User `json:"items"`
}

type UserChangePassword struct {
	Name     string `json:"name" validate:"required"`
	Password string `json:"password" validate:"kopassword,required,max=30,min=8"`
	Original string `json:"original" validate:"kopassword,required,max=30,min=8"`
}

type UserForgotPassword struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
}
