package v1

import (
	"github.com/gin-gonic/gin"
	"ko3-gin/pkg/router/v1/host"
)

func V1(root *gin.RouterGroup) *gin.RouterGroup {
	v1Api := root.Group("v1")
	{
		v1HostApi := v1Api.Group("/hosts")
		{
			v1HostApi.GET("/", host.List)
			v1HostApi.POST("/", host.Create)
			v1HostApi.GET("/:name/", host.Get)
			v1HostApi.PATCH("/:name/", host.Update)
			v1HostApi.DELETE("/:name/", host.Delete)
			v1HostApi.POST("/batch/:name", host.Batch)
		}
	}
	return v1Api
}
