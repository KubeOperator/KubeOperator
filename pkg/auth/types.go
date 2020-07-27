package auth

import "github.com/KubeOperator/KubeOperator/pkg/permission"

type Credential struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Language string `json:"language"`
}

type SessionUser struct {
	UserId   string `json:"userId"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Language string `json:"language"`
	IsActive bool   `json:"isActive"`
	IsAdmin  bool   `json:"isAdmin"`
}

type JwtResponse struct {
	User        SessionUser                 `json:"user"`
	Token       string                      `json:"token"`
	RoleMenus   []permission.UserMenu       `json:"roleMenus"`
	Permissions []permission.UserPermission `json:"permissions"`
}
