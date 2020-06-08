package proxy

import (
	"crypto/tls"
	"fmt"
	clusterService "github.com/KubeOperator/KubeOperator/pkg/service/cluster"
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
	c, err := clusterService.Get(clusterName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	endpoint, err := clusterService.GetDefaultClusterEndpoint(clusterName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	u, err := url.Parse(fmt.Sprintf("http://%s", endpoint))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	proxy := httputil.NewSingleHostReverseProxy(u)
	proxy.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	ctx.Request.Host = fmt.Sprintf("logging.%s", c.Spec.AppDomain)
	ctx.Request.URL.Path = path
	proxy.ServeHTTP(ctx.Writer, ctx.Request)
}
