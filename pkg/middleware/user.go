package middleware

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/kataras/iris/v12/context"
	"net/http"
)

func UserMiddleware(ctx context.Context) {
	//user := ctx.Values().Get("jwt").(*jwt.Token)
	//foobar := user.Claims.(jwt.MapClaims)
	//sessionUserJson, _ := json.Marshal(foobar)
	//sessionUserJsonStr := string(sessionUserJson)
	session := constant.Sess.Start(ctx)
	sessionUser := session.Get(constant.SessionUserKey)
	if sessionUser == nil {
		ctx.StatusCode(http.StatusUnauthorized)
		ctx.StopExecution()
		return
	}
	user := sessionUser.(*dto.Profile)
	ctx.Values().Set("user", user)
	ctx.Next()
}
