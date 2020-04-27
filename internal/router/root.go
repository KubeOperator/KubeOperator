package router

import (
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	_ "ko3-gin/docs"
	"ko3-gin/internal/middleware"
	pkg_api "ko3-gin/pkg/api"
	"os"
)


func Server() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	f, _ := os.Open(os.DevNull)
	gin.DefaultWriter = f
	server := gin.Default()
	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	server.Use(middleware.LoggerMiddleware())
	jwtMiddleware := middleware.JWTMiddleware()
	auth := server.Group("/auth")
	{
		auth.POST("/login", jwtMiddleware.LoginHandler)
		auth.GET("/refresh", jwtMiddleware.RefreshHandler)
	}
	api := server.Group("/api")
	api.Use(jwtMiddleware.MiddlewareFunc())
	{
		pkg_api.V1(api)
	}
	return server
}
