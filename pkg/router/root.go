package router

import (
	"github.com/KubeOperator/KubeOperator/pkg/middleware"
	v1 "github.com/KubeOperator/KubeOperator/pkg/router/v1"
	"github.com/kataras/iris/v12"
)

func Server() *iris.Application {
	app := iris.New()
	_ = app.I18n.Load("../pkg/locales/*/*", "en-US", "zh-CN")
	app.I18n.SetDefault("zh-CN")
	app.Post("/api/auth/login", middleware.LoginHandler)
	app.Use(middleware.LogMiddleware)
	app.Use(middleware.JWTMiddleware().Serve)
	api := app.Party("/api")
	v1.V1(api)
	return app
}
