package auth

import (
	"errors"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/service"
	"github.com/dgrijalva/jwt-go"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/mojocn/base64Captcha"
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

	if aul.CaptchaId != "" {
		err := VerifyCode(aul.CaptchaId, aul.Code)
		if err != nil {
			ctx.StatusCode(iris.StatusUnauthorized)
			response := new(dto.Response)
			if ctx.Tr(err.Error()) == "" {
				response.Msg = err.Error()
			} else {
				response.Msg = ctx.Tr(err.Error())
			}
			_, _ = ctx.JSON(response)
			return
		}
	}

	data, err := CheckLogin(aul.Username, aul.Password)
	if err != nil {
		ctx.StatusCode(iris.StatusUnauthorized)
		response := new(dto.Response)
		if ctx.Tr(err.Error()) == "" {
			response.Msg = err.Error()
		} else {
			response.Msg = ctx.Tr(err.Error())
		}
		_, _ = ctx.JSON(response)
		return
	}
	ctx.StatusCode(iris.StatusOK)
	_, _ = ctx.JSON(data)
	return
}

var store = base64Captcha.DefaultMemStore
var verifyCodeFailed = errors.New("VERIFY_CODE_FAILED")

func VerificationCodeHandler(ctx context.Context) {
	var driverString base64Captcha.DriverString
	driverString.Source = "1234567890qwertyuioplkjhgfdsazxcvbnm"
	driverString.Width = 120
	driverString.Height = 50
	driverString.NoiseCount = 0
	driverString.Length = 4
	driverString.Fonts = []string{"wqy-microhei.ttc"}
	driver := driverString.ConvertFonts()
	c := base64Captcha.NewCaptcha(driver, store)
	id, b64s, err := c.Generate()
	body := map[string]interface{}{"code": 1, "image": b64s, "captchaId": id, "msg": "success"}
	if err != nil {
		body = map[string]interface{}{"code": 0, "error": map[string]interface{}{"msg": err.Error()}}
	}
	ctx.Header("Content-Type", "application/json; charset=utf-8")
	ctx.JSON(body)
}

func VerifyCode(codeId string, code string) error {
	if code == "" {
		return verifyCodeFailed
	}
	if store.Verify(codeId, code, true) {
		return nil
	} else {
		return verifyCodeFailed
	}
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
