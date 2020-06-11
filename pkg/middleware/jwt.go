package middleware

import (
	"github.com/KubeOperator/KubeOperator/pkg/auth"
	"github.com/KubeOperator/KubeOperator/pkg/logger"
	"github.com/KubeOperator/KubeOperator/pkg/service"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
	"time"
)

var log = logger.Default
var sessionUser = auth.SessionUser{}

func JWTMiddleware() *jwt.GinJWTMiddleware {
	secret := viper.GetString("app.secret")
	j, err := jwt.New(&jwt.GinJWTMiddleware{
		Key:           []byte(secret),
		Timeout:       time.Hour,
		MaxRefresh:    time.Hour,
		TimeFunc:      time.Now,
		TokenHeadName: "Bearer",
		IdentityKey:   "user",
		Authenticator: func(ctx *gin.Context) (i interface{}, err error) {
			var credential auth.Credential
			if err := ctx.ShouldBind(&credential); err != nil {
				return nil, jwt.ErrMissingLoginValues
			}
			sUser, err := service.UserAuth(credential.Username, credential.Password)
			if err != nil {
				if sUser != nil && sUser.IsActive == false {
					return nil, err
				} else {
					return nil, err
				}
			}
			return sUser, nil
		},
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*auth.SessionUser); ok {
				sessionUser = *v
				return jwt.MapClaims{
					"user": v,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(ctx *gin.Context) interface{} {
			claims := jwt.ExtractClaims(ctx)
			return &claims
		},
		LoginResponse: func(ctx *gin.Context, code int, token string, expire time.Time) {
			ctx.JSON(http.StatusOK, gin.H{
				"code":   code,
				"token":  token,
				"expire": expire.Format(time.RFC3339),
				"user":   sessionUser,
			})
			return
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	return j
}

func GetAuthUser(ctx *gin.Context) {
	claims := jwt.ExtractClaims(ctx)
	token := ctx.Keys["JWT_TOKEN"].(string)
	ctx.JSON(http.StatusOK, gin.H{
		"code":  http.StatusOK,
		"token": token,
		"user":  claims["user"],
	})
	return
}
