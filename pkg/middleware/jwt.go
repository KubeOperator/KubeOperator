package middleware

import (
	"errors"
	"net/http"
	"time"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/dgrijalva/jwt-go"
	jwtmiddleware "github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris/v12/context"
	"github.com/spf13/viper"
)

var (
	UserIsNotRelatedProject = "USER_IS_NOT_RELATED_PROJECT"
)

type JwtMiddleware struct {
	*jwtmiddleware.Middleware
}

func (m *JwtMiddleware) Serve(ctx context.Context) {
	session := constant.Sess.Start(ctx)
	u := session.Get(constant.SessionUserKey)
	user, ok := u.(*dto.Profile)
	if ok {
		if user.Timeout.Before(time.Now()) {
			m.Config.ErrorHandler(ctx, errors.New("token timeout !"))
			session.Delete(constant.SessionUserKey)
			return
		}

		exp := viper.GetInt("jwt.exp")
		user.Timeout = time.Now().Add((time.Minute * time.Duration(exp)))
		session.Set(constant.SessionUserKey, user)
		ctx.Next()
		return
	}
	ctx.Next()
}

func JWTMiddleware() *JwtMiddleware {
	secretKey := []byte(viper.GetString("jwt.secret"))
	m := JwtMiddleware{jwtmiddleware.New(
		jwtmiddleware.Config{
			Extractor: jwtmiddleware.FromAuthHeader,
			ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
				return secretKey, nil
			},
			SigningMethod: jwt.SigningMethodHS256,
			ErrorHandler:  ErrorHandler,
		},
	)}
	return &m
}

func ErrorHandler(ctx context.Context, err error) {
	if err == nil {
		return
	}
	ctx.StopExecution()
	response := &dto.Response{
		Msg: err.Error(),
	}
	ctx.StatusCode(http.StatusUnauthorized)
	_, _ = ctx.JSON(response)
}
