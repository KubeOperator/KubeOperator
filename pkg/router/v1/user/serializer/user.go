package serializer

import (
	userModel "github.com/KubeOperator/KubeOperator/pkg/model/user"
	"time"
)

type User struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	Password string    `json:"password"`
	Email    string    `json:"email"`
	IsActive bool      `json:"isActive"`
	Language string    `json:"language"`
	CreateAt time.Time `json:"createAt"`
	UpdateAt time.Time `json:"updateAt"`
}

func FromModel(u userModel.User) User {
	return User{
		ID:       u.ID,
		Name:     u.Name,
		Password: u.Password,
		Email:    u.Email,
		IsActive: u.IsActive,
		Language: u.Language,
		CreateAt: u.CreatedAt,
		UpdateAt: u.UpdatedAt,
	}
}

func ToModel(u User) userModel.User {
	return userModel.User{
		ID:       u.ID,
		Name:     u.Name,
		Password: u.Password,
		Email:    u.Email,
		Language: u.Language,
	}
}

type ListUserResponse struct {
	Items []User `json:"items"`
	Total int    `json:"total"`
}

type GetUserResponse struct {
	Item User `json:"item"`
}

type CreateUserRequest struct {
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required"`
}

type CreateUserResponse struct {
	Item User `json:"item"`
}

type DeleteUserRequest struct {
	Name string `json:"name"`
}

type DeleteUserResponse struct {
}

type UpdateUserRequest struct {
	ID       string `json:"id" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Password string `json:"password"`
	Email    string `json:"email"`
	IsActive bool   `json:"isActive"`
	Language string `json:"language"`
}

type UpdateUserResponse struct {
	Item User `json:"item"`
}

type BatchUserRequest struct {
	Operation string `json:"operation" binding:"required"`
	Items     []User `json:"items"`
}

type BatchUserResponse struct {
	Items []User `json:"items"`
}

type ChangePasswordRequest struct {
	Name     string `json:"name" binding:"required"`
	Original string `json:"original"`
	Password string `json:"password" binding:"required"`
}
