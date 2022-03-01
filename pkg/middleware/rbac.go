package middleware

import (
	"net/http"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/logger"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/storyicon/grbac"
)

var log = logger.Default

func RBACMiddleware() iris.Handler {
	r, err := grbac.New(grbac.WithAdvancedRules(constant.SystemRules))
	if err != nil {
		panic(err)
	}
	return func(c context.Context) {
		roles := querySystemRoles(c)
		state, err := r.IsRequestGranted(c.Request(), roles)
		if err != nil {
			c.StatusCode(http.StatusInternalServerError)
			c.StopExecution()
			return
		}
		if !state.IsGranted() {
			c.StatusCode(http.StatusForbidden)
			c.StopExecution()
			return
		}
		c.Next()
	}
}

func querySystemRoles(ctx context.Context) []string {
	user := ctx.Values().Get("user")
	sessionUser, ok := user.(dto.SessionUser)
	if !ok {
		log.Errorf("type aassertion failed")
	}
	roles := sessionUser.Roles
	ctx.Values().Set("roles", roles)
	return roles
}
