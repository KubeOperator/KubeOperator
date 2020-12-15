package proxy

import (
	"crypto/tls"
	"fmt"
	kubeUtil "github.com/KubeOperator/KubeOperator/pkg/util/kubernetes"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func KubernetesClientProxy(ctx context.Context) {
	clusterName := ctx.Params().Get("cluster_name")
	proxyPath := ctx.Params().Get("p")
	endpoints, err := clusterService.GetApiServerEndpoints(clusterName)

	aliveHost, err := kubeUtil.SelectAliveHost(endpoints)
	if err != nil {
		_, _ = ctx.JSON(iris.StatusInternalServerError)
		return
	}
	u, err := url.Parse(fmt.Sprintf("https://%s", aliveHost))
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
	proxy.ModifyResponse = func(response *http.Response) error {
		if response.StatusCode == http.StatusUnauthorized {
			response.StatusCode = http.StatusInternalServerError
		}
		return nil
	}
	proxy.ServeHTTP(ctx.ResponseWriter(), ctx.Request())
}
