package proxy

import (
	"crypto/tls"
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"net/http"
	"net/http/httputil"
	"net/url"
)



func KubernetesClientProxy(ctx context.Context) {
	clusterName := ctx.Params().Get("cluster_name")
	proxyPath := ctx.Params().Get("p")
	api, err := clusterService.GetEndpoint(clusterName)
	if err != nil {
		_, _ = ctx.JSON(iris.StatusInternalServerError)
		return
	}
	u, err := url.Parse(fmt.Sprintf("https://%s:8443", api))
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
