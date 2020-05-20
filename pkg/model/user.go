package model

import "github.com/KubeOperator/KubeOperator/pkg/auth"

type User struct {
	Id       string `gorm:"primary_key size:64" `
	Name     string `gorm:"size:128"`
	Password string `gorm:"size:256"`
}

func (u *User)ToSessionUser() *auth.SessionUser  {
	return &auth.SessionUser{
		UserId: u.Id,
		Name:   u.Name,
	}
}
