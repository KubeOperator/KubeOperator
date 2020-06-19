package proxy

import (
	"crypto/tls"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func PrometheusProxy(ctx context.Context) {
	proxyPath := ctx.Params().Get("p")
	u, err := url.Parse("http://localhost:3000")
	if err != nil {
		_, _ = ctx.JSON(iris.StatusInternalServerError)
		return
	}
	proxy := httputil.NewSingleHostReverseProxy(u)
	proxy.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	ctx.Request().URL.Path = proxyPath
	proxy.ServeHTTP(ctx.ResponseWriter(), ctx.Request())
}
