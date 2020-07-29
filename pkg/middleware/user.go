package middleware

import (
	"encoding/json"
	"github.com/KubeOperator/KubeOperator/pkg/auth"
	"github.com/dgrijalva/jwt-go"
	"github.com/kataras/iris/v12/context"
)

func UserMiddleware(ctx context.Context) {
	user := ctx.Values().Get("jwt").(*jwt.Token)
	foobar := user.Claims.(jwt.MapClaims)
	sessionUserJson, _ := json.Marshal(foobar)
	sessionUserJsonStr := string(sessionUserJson)
	var sessionUser auth.SessionUser
	json.Unmarshal([]byte(sessionUserJsonStr), &sessionUser)
	ctx.Values().Set("user", sessionUser)
	ctx.Next()
}
