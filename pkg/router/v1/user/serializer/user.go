package serializer

import (
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	userModel "github.com/KubeOperator/KubeOperator/pkg/model/user"
)

type User struct {
	common.BaseModel
	ID       string `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Email    string `json:"email"`
	IsActive bool   `json:"isActive"`
	Language string `json:"language"`
}

func FromModel(u userModel.User) User {
	return User{
		ID:       u.ID,
		Name:     u.Name,
		Password: u.Password,
		Email:    u.Email,
		IsActive: u.IsActive,
		Language: u.Language,
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
