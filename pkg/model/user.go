package model

import (
	"github.com/KubeOperator/KubeOperator/pkg/auth"
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	uuid "github.com/satori/go.uuid"
)

const (
	EN string = "en"
	ZH string = "zh"
)

type User struct {
	common.BaseModel
	ID       string
	Name     string `gorm:"not null;unique"`
	Password string
	Email    string `gorm:"not null;unique"`
	IsActive bool
	Language string
}

type Token struct {
	Token string `json:"access_token"`
}

func (u *User) BeforeCreate() (err error) {
	u.ID = uuid.NewV4().String()
	return err
}

func (u User) TableName() string {
	return "ko_user"
}

func (u *User) ToSessionUser() *auth.SessionUser {
	return &auth.SessionUser{
		UserId:   u.ID,
		Name:     u.Name,
		Email:    u.Email,
		Language: u.Language,
		IsActive: u.IsActive,
	}
}
