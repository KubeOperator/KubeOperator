package proxy

import (
	"crypto/tls"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func LoggingProxy(ctx *gin.Context) {
	clusterName := ctx.Param("name")
	path := ctx.Param("path")
	if clusterName == "" {
		ctx.JSON(http.StatusBadRequest, invalidClusterNameError)
		return
	}
	u, err := url.Parse("http://172.16.10.184")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	proxy := httputil.NewSingleHostReverseProxy(u)
	proxy.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	ctx.Request.Host = "logging.adfsfssf.com"
	ctx.Request.URL.Path = path
	proxy.ServeHTTP(ctx.Writer, ctx.Request)
}
