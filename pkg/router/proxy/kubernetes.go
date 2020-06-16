package proxy

import (
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/KubeOperator/KubeOperator/pkg/service"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"net/http"
	"net/http/httputil"
	"net/url"
)

var (
	invalidClusterNameError = errors.New("invalid cluster name")
	keyPrefix               = "Bearer"
	AuthorizationHeader     = "Authorization"
)

func KubernetesClientProxy(ctx context.Context) {
	clusterName := ctx.URLParam("cluster_name")
	proxyPath := ctx.URLParam("p")
	var clusterService service.ClusterService
	api, err := clusterService.GetEndpoint(clusterName)
	if err != nil {
		_, _ = ctx.JSON(iris.StatusInternalServerError)
		return
	}
	u, err := url.Parse(api)
	if err != nil {
		_, _ = ctx.JSON(iris.StatusInternalServerError)
		return
	}
	proxy := httputil.NewSingleHostReverseProxy(u)
	proxy.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	secret, err := clusterService.GetSecrets(clusterName)
	if err != nil {
		_, _ = ctx.JSON(iris.StatusInternalServerError)
		return
	}
	token := fmt.Sprintf("%s %s", keyPrefix, secret.KubernetesToken)
	ctx.Request().Header.Add(AuthorizationHeader, token)
	ctx.Request().URL.Path = proxyPath
	proxy.ServeHTTP(ctx.ResponseWriter(), ctx.Request())
}
