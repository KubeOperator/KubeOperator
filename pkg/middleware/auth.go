package middleware

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	jwtmiddleware "github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris/v12/context"
)

type AuthMiddleware struct {
	*jwtmiddleware.Middleware
}

func (m *AuthMiddleware) Serve(ctx context.Context) {
	session := constant.Sess.Start(ctx)
	u := session.Get(constant.SessionUserKey)
	if u != nil {
		ctx.Next()
		return
	}
	if err := m.CheckJWT(ctx); err != nil {
		m.Config.ErrorHandler(ctx, err)
		return
	}
	ctx.Next()
}
