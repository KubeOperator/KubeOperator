package router

import (
	"github.com/KubeOperator/KubeOperator/pkg/controller"
	"github.com/KubeOperator/KubeOperator/pkg/middleware"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
)

func Server() *iris.Application {
	app := iris.New()
	app.Use(middleware.LogMiddleware)
	api := app.Party("/api")
	v1 := api.Party("/v1")
	mvc.New(v1.Party("/demo")).Handle(controller.NewDemoController())
	return app
}
