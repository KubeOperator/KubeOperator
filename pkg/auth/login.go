package auth

import (
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/service"
	"github.com/dgrijalva/jwt-go"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/spf13/viper"
	"time"
)

// Login
// @Tags auth
// @Summary Get login token
// @Description Get login token
// @Param request body auth.Credential true "request"
// @Accept  json
// @Produce  json
// @Success 200 {object} auth.JwtResponse
// @Router /auth/login/ [post]
func LoginHandler(ctx context.Context) {
	aul := new(Credential)
	if err := ctx.ReadJSON(&aul); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		_, _ = ctx.JSON(dto.Response{Msg: err.Error()})
		return
	}

	data, err := CheckLogin(aul.Username, aul.Password)
	if err != nil {
		ctx.StatusCode(iris.StatusUnauthorized)
		_, _ = ctx.JSON(dto.Response{Msg: ctx.Tr(err.Error())})
		return
	}
	ctx.StatusCode(iris.StatusOK)
	_, _ = ctx.JSON(data)
	return
}
func CheckLogin(username string, password string) (*JwtResponse, error) {
	user, err := service.UserAuth(username, password)
	if err != nil {
		return nil, err
	}
	token, err := CreateToken(toSessionUser(*user))
	if err != nil {
		return nil, err
	}
	resp := new(JwtResponse)
	resp.Token = token
	resp.User = toSessionUser(*user)
	return resp, err
}

func CreateToken(user SessionUser) (string, error) {
	exp := viper.GetInt("jwt.exp")
	secretKey := []byte(viper.GetString("jwt.secret"))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"name":     user.Name,
		"email":    user.Email,
		"userId":   user.UserId,
		"isActive": user.IsActive,
		"language": user.Language,
		"isAdmin":  user.IsAdmin,
		"iat":      time.Now().Unix(),
		"exp":      time.Now().Add(time.Minute * time.Duration(exp)).Unix(),
	})
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return tokenString, err
}

func toSessionUser(u model.User) SessionUser {
	return SessionUser{
		UserId:   u.ID,
		Name:     u.Name,
		Email:    u.Email,
		Language: u.Language,
		IsActive: u.IsActive,
		IsAdmin:  u.IsAdmin,
	}
}
