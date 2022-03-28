package proxy

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/kataras/iris/v12/context"
)

func LoggingProxy(ctx context.Context) {
	clusterName := ctx.Params().Get("cluster_name")
	proxyPath := ctx.Params().Get("p")
	if clusterName == "" {
		_, _ = ctx.JSON(http.StatusBadRequest)
		return
	}

	nodeInfo, err := clusterToolService.GetNodePort(clusterName, "loki")
	if err != nil {
		_, _ = ctx.JSON(http.StatusInternalServerError)
		return
	}

	u, err := url.Parse(fmt.Sprintf("http://%s:%v", nodeInfo.NodeHost, nodeInfo.NodePort))
	if err != nil {
		_, _ = ctx.JSON(http.StatusInternalServerError)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(u)
	ctx.Request().URL.Path = proxyPath
	proxy.ServeHTTP(ctx.ResponseWriter(), ctx.Request())
}
