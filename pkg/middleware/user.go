package middleware

import (
	"encoding/json"
	"net/http"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/dgrijalva/jwt-go"
	"github.com/kataras/iris/v12/context"
)

func UserMiddleware(ctx context.Context) {
	var u dto.SessionUser
	j := ctx.Values().Get("jwt")
	if j != nil {
		j := j.(*jwt.Token)
		foobar := j.Claims.(jwt.MapClaims)
		js, _ := json.Marshal(foobar)
		_ = json.Unmarshal(js, &u)
	} else {
		session := constant.Sess.Start(ctx)
		sessionUser := session.Get(constant.SessionUserKey)
		if sessionUser == nil {
			ctx.StatusCode(http.StatusUnauthorized)
			ctx.StopExecution()
			return
		}
		u = sessionUser.(*dto.Profile).User
	}
	// set roles
	ctx.Values().Set("user", u)
	ctx.Values().Set("operator", u.Name)
	ctx.Next()
}
