package proxy

import (
	"crypto/tls"
	"fmt"
	"github.com/KubeOperator/KubeOperator/pkg/service"
	"github.com/kataras/iris/v12/context"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func LoggingProxy(ctx context.Context) {
	var clusterService service.ClusterService
	clusterName := ctx.URLParam("cluster_name")
	proxyPath := ctx.URLParam("p")
	if clusterName == "" {
		_, _ = ctx.JSON(http.StatusBadRequest)
		return
	}
	c, err := clusterService.Get(clusterName)
	if err != nil {
		_, _ = ctx.JSON(http.StatusInternalServerError)
		return
	}

	endpoint, err := clusterService.GetEndpoint(clusterName)
	if err != nil {
		_, _ = ctx.JSON(http.StatusInternalServerError)
		return
	}
	u, err := url.Parse(fmt.Sprintf("http://%s", endpoint))
	if err != nil {
		_, _ = ctx.JSON(http.StatusInternalServerError)
		return
	}
	proxy := httputil.NewSingleHostReverseProxy(u)
	proxy.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	ctx.Request().Host = fmt.Sprintf("logging.%s", c.Spec.AppDomain)
	ctx.Request().URL.Path = proxyPath
	proxy.ServeHTTP(ctx.ResponseWriter(), ctx.Request())
}
