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
