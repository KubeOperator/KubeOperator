package session

import (
	"errors"
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/kataras/iris/v12/context"
)

func GetUser(ctx context.Context) (*dto.Profile, error) {
	session := constant.Sess.Start(ctx)
	user := session.Get(constant.SessionUserKey)
	if user == nil {
		return nil, errors.New("user is not login")
	}
	p, ok := user.(*dto.Profile)
	if !ok {
		return nil, errors.New("can not parse to user profile")
	}
	return p, nil
}

func GetProjectName(ctx context.Context) (string, error) {
	session := constant.Sess.Start(ctx)
	sessionUser := session.Get(constant.SessionUserKey)
	if sessionUser == nil {
		return "", errors.New("user is not login")
	}
	projectName := ""
	profile, ok := sessionUser.(*dto.Profile)
	if !ok {
		return "", errors.New("can not parse to user profile")
	}
	user := profile.User
	if !user.IsAdmin {
		projectName = user.CurrentProject
	}
	return projectName, nil
}
