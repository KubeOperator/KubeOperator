package user

import (
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
)

type User struct {
	common.BaseModel
	ID       string
	Name     string
	Password string
	Email    string
	IsActive bool
	Language string
}

func (u User) TableName() string {
	return "ko_user"
}
