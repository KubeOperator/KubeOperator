package middleware

import (
	"errors"
	"net/http"
	"time"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/session"
	"github.com/kataras/iris/v12/context"
)

func SessionMiddleware(ctx context.Context) {
	var sessionID = session.GloablSessionMgr.CheckCookieValid(ctx.ResponseWriter(), ctx.Request())
	if sessionID == "" {
		errorHandler(ctx, errors.New("session invalid !"))
		return
	}

	u, ok := session.GloablSessionMgr.GetSessionVal(sessionID, constant.SessionUserKey)
	if !ok {
		errorHandler(ctx, errors.New("session invalid !"))
		return
	}

	user, ok := u.(*dto.Profile)
	if ok {
		if session.GloablSessionMgr.GetLastAccessTime(sessionID).Add(time.Second * 3600).Before(time.Now()) {
			session.GloablSessionMgr.EndSessionBy(sessionID)
			errorHandler(ctx, errors.New("token timeout !"))
			return
		}

		ctx.Values().Set("user", user.User)
		ctx.Values().Set("operator", user.User.Name)
		ctx.Next()
		return
	}
	ctx.Next()
}

func errorHandler(ctx context.Context, err error) {
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
