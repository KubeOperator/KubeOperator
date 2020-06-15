package middleware

import (
	"github.com/kataras/iris/v12/context"
)

func LogMiddleware(ctx context.Context) {
	ctx.Application().Logger().Infof("Path: %s | IP: %s", ctx.Path(), ctx.RemoteAddr())
	ctx.Next()
}
