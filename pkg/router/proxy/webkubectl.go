package proxy

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"net/http/httputil"
	"net/url"
)

func WebKubeCtlProxy(ctx context.Context) {
	proxyPath := ctx.Params().Get("p")
	u, err := url.Parse("http://localhost:8082")
	if err != nil {
		_, _ = ctx.JSON(iris.StatusInternalServerError)
		return
	}
	proxy := httputil.NewSingleHostReverseProxy(u)
	ctx.Request().URL.Path = proxyPath
	proxy.ServeHTTP(ctx.ResponseWriter(), ctx.Request())
}
