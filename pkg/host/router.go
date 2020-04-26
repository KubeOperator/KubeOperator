package host

import "github.com/gin-gonic/gin"

func Router(parent *gin.RouterGroup) {
	g := parent.Group("/host")
	{
		g.POST("/", Create)
	}
}
