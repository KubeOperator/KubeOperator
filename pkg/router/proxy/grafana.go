package proxy

import (
	"crypto/tls"
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/spf13/viper"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func GrafanaProxy(ctx context.Context) {
	proxyPath := ctx.Params().Get("p")
	grafanaUrl := fmt.Sprintf("http://%s:%d", viper.GetString("grafana.host"), viper.GetInt("grafana.port"))
	u, err := url.Parse(grafanaUrl)
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
