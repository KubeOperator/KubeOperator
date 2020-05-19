package v1

import (
	"github.com/gin-gonic/gin"
	"ko3-gin/pkg/router/v1/cluster"
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
			v1HostApi.POST("/batch/", host.Batch)
		}
		v1ClusterApi := v1Api.Group("/clusters")
		{
			v1ClusterApi.GET("/", cluster.List)
			v1ClusterApi.POST("/", cluster.Create)
			v1ClusterApi.GET("/:name/", cluster.Get)
			v1ClusterApi.PATCH("/:name/", cluster.Update)
			v1ClusterApi.DELETE("/:name/", cluster.Delete)
			v1ClusterApi.POST("/batch/", cluster.Batch)
			v1ClusterApi.POST("/page/", cluster.Page)
		}
	}
	return v1Api
}
