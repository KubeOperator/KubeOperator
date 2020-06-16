package router

import (
	"github.com/KubeOperator/KubeOperator/pkg/middleware"
	"github.com/KubeOperator/KubeOperator/pkg/router/proxy"
	v1 "github.com/KubeOperator/KubeOperator/pkg/router/v1"
	"github.com/kataras/iris/v12"
)

func Server() *iris.Application {
	app := iris.New()
	_ = app.I18n.Load("../pkg/locales/*/*", "en-US", "zh-CN")
	app.I18n.SetDefault("zh-CN")
	app.Post("/api/auth/login", middleware.LoginHandler)
	proxy.RegisterProxy(app)
	api := app.Party("/api")
	api.Use(middleware.PagerMiddleware)
	//api.Use(middleware.LogMiddleware)
	//api.Use(middleware.JWTMiddleware().Serve)
	v1.V1(api)
	return app
}
