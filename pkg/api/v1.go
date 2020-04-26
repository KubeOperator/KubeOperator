package api

import "github.com/gin-gonic/gin"

func V1(root *gin.RouterGroup) *gin.RouterGroup {
	router := root.Group("v1")
	return router
}
