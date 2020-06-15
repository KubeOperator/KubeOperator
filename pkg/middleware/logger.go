package middleware

import "github.com/kataras/iris"

func LogMiddleware(ctx iris.Context) {
	ctx.Application().Logger().Infof("Path: %s | IP: %s", ctx.Path(), ctx.RemoteAddr())
	ctx.Next()
}

