package service

import (
	"errors"
	"github.com/KubeOperator/KubeOperator/pkg/auth"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/i18n"
	"github.com/KubeOperator/KubeOperator/pkg/model/user"
	"github.com/KubeOperator/KubeOperator/pkg/util/encrypt"
)

var (
	UserNotFound     = errors.New(i18n.Tr("user_not_found", nil))
	PasswordNotMatch = errors.New(i18n.Tr("password_not_match", nil))
	UserIsNotActive  = errors.New(i18n.Tr("user_is_not_active", nil))
)

func UserAuth(name string, password string) (sessionUser *auth.SessionUser, err error) {
	var dbUser user.User
	if db.DB.Where("name = ?", name).First(&dbUser).RecordNotFound() {
		return nil, UserNotFound
	}
	if dbUser.IsActive == false {
		return dbUser.ToSessionUser(), UserIsNotActive
	}
	password, err = encrypt.StringEncrypt(password)
	if err != nil {
		return nil, err
	}
	if dbUser.Password != password {
		return nil, PasswordNotMatch
	}
	return dbUser.ToSessionUser(), nil
}
