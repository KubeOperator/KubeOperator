package proxy

import (
	"crypto/tls"
	"errors"
	"fmt"
	clusterService "github.com/KubeOperator/KubeOperator/pkg/service/cluster"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httputil"
	"net/url"
)

var (
	invalidClusterNameError = errors.New("invalid cluster name")
	keyPrefix               = "Bearer"
	AuthorizationHeader     = "Authorization"
)

func KubernetesClientProxy(ctx *gin.Context) {
	clusterName := ctx.Param("name")
	path := ctx.Param("path")
	if clusterName == "" {
		ctx.JSON(http.StatusBadRequest, invalidClusterNameError)
		return
	}
	api, err := clusterService.GetClusterKubernetesApiEndpoint(clusterName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	u, err := url.Parse(api)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	proxy := httputil.NewSingleHostReverseProxy(u)
	proxy.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	secret, err := clusterService.GetClusterSecret(clusterName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	token := fmt.Sprintf("%s %s", keyPrefix, secret.KubernetesToken)
	ctx.Request.Header.Add(AuthorizationHeader, token)
	ctx.Request.URL.Path = path
	proxy.ServeHTTP(ctx.Writer, ctx.Request)

}
