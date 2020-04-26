package api

import (
	"github.com/gin-gonic/gin"
	"ko3-gin/pkg/host"
)

func V1(root *gin.RouterGroup) *gin.RouterGroup {
	router := root.Group("v1")
	{
		host.Router(router)
	}
	return router
}
