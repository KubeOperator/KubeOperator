package router

import (
	_ "github.com/KubeOperator/KubeOperator/docs"
	"github.com/KubeOperator/KubeOperator/pkg/middleware"
	v1 "github.com/KubeOperator/KubeOperator/pkg/router/v1"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"net/http"
	"os"
)

func Server() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	f, _ := os.Open(os.DevNull)
	gin.DefaultWriter = f
	server := gin.Default()
	server.StaticFS("static", http.Dir("resource/static"))
	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	server.NoRoute(NotFoundResponse)
	server.Use(middleware.LoggerMiddleware())
	server.Use(middleware.PagerMiddleware())
	api := server.Group("/api")
	{
		v1.V1(api)
	}

	jwtMiddleware := middleware.JWTMiddleware()
	auth := server.Group("/auth")
	{
		auth.POST("/login", jwtMiddleware.LoginHandler)
		auth.GET("/refresh", jwtMiddleware.RefreshHandler)
	}
	auth.Use(jwtMiddleware.MiddlewareFunc())
	//api := server.Group("/api")
	//api.Use(jwtMiddleware.MiddlewareFunc())
	//{
	//	PkgApi.V1(api)
	//}
	return server
}

func NotFoundResponse(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{
		"status": 404,
		"error":  "not found",
	})
}
