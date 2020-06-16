package controller

import (
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/KubeOperator/KubeOperator/pkg/service"
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

func NewProxyController() *proxyController {
	return &proxyController{}
}

type proxyController struct {
	Ctx            context.Context
	clusterService service.ClusterService
}

func (p proxyController) AnyKubernetes(clusterName string, path string) error {
	api, err := p.clusterService.GetEndpoint(clusterName)
	if err != nil {
		return err
	}
	u, err := url.Parse(api)
	if err != nil {
		return err
	}
	proxy := httputil.NewSingleHostReverseProxy(u)
	proxy.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	secret, err := p.clusterService.GetSecrets(clusterName)
	if err != nil {
		return err
	}
	token := fmt.Sprintf("%s %s", keyPrefix, secret.KubernetesToken)
	p.Ctx.Request().Header.Add(AuthorizationHeader, token)
	p.Ctx.Request().URL.Path = path
	proxy.ServeHTTP(p.Ctx.ResponseWriter(), p.Ctx.Request())
	return nil
}
