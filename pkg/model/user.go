package model

import (
	"github.com/KubeOperator/KubeOperator/pkg/auth"
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	uuid "github.com/satori/go.uuid"
)

const (
	EN string = "en-US"
	ZH string = "zh-CN"
)

type User struct {
	common.BaseModel
	ID       string `json:"id" gorm:"type:varchar(64)"`
	Name     string `json:"name" gorm:"type:varchar(256);not null;unique"`
	Password string `json:"password" gorm:"type:varchar(256)"`
	Email    string `json:"email" gorm:"type:varchar(256);not null;unique"`
	IsActive bool   `json:"isActive" gorm:"type:int(1)"`
	Language string `json:"language" gorm:"type:varchar(64)"`
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
