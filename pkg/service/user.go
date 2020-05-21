package service

import (
	"errors"
	"github.com/KubeOperator/KubeOperator/pkg/auth"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"golang.org/x/crypto/bcrypt"
)

var (
	UserNotFound     = errors.New("can not find user")
	PasswordNotMatch = errors.New("password not match")
)

func UserAuth(name string, password string) (sessionUser *auth.SessionUser, err error) {
	var user model.User
	if db.DB.Where("name = ?", name).First(&user).RecordNotFound() {
		return nil, UserNotFound
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, PasswordNotMatch
	}
	return user.ToSessionUser(), nil
}
