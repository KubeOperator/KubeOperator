package proxy

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/kataras/iris/v12/context"
)

func GrafanaProxy(ctx context.Context) {
	clusterName := ctx.Params().Get("cluster_name")
	proxyPath := ctx.Params().Get("p")
	if clusterName == "" {
		_, _ = ctx.JSON(http.StatusBadRequest)
		return
	}
	endpoint, err := clusterService.GetRouterEndpoint(clusterName)
	if err != nil {
		_, _ = ctx.JSON(http.StatusInternalServerError)
		return
	}
	u, err := url.Parse(fmt.Sprintf("http://%s", endpoint.Address))
	if err != nil {
		_, _ = ctx.JSON(http.StatusInternalServerError)
		return
	}
	proxy := httputil.NewSingleHostReverseProxy(u)
	proxy.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	req := ctx.Request()
	req.Host = fmt.Sprintf(constant.DefaultGrafanaIngress)
	if proxyPath == "root" {
		proxyPath = "/"
	}
	req.URL.Path = proxyPath
	proxy.ServeHTTP(ctx.ResponseWriter(), req)
}
