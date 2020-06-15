package router

import (
	"github.com/KubeOperator/KubeOperator/pkg/controller"
	"github.com/KubeOperator/KubeOperator/pkg/middleware"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

func Server() *iris.Application {
	app := iris.New()
	app.Post("/api/auth/login", middleware.LoginHandler)
	app.Use(middleware.LogMiddleware)
	app.Use(middleware.JWTMiddleware().Serve)
	api := app.Party("/api")
	v1 := api.Party("/v1")
	mvc.New(v1.Party("/demo")).Handle(controller.NewDemoController())
	return app
}
