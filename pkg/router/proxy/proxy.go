package proxy

import "github.com/gin-gonic/gin"

func Proxy(proxyRoot *gin.RouterGroup) *gin.RouterGroup {
	proxyRoot.GET("/kubernetes/:name/*path", KubernetesClientProxy)
	proxyRoot.GET("/logging/:name/*path", LoggingProxy)
	proxyRoot.POST("/logging/:name/*path", LoggingProxy)
	return proxyRoot
}
