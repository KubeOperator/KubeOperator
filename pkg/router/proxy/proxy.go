package proxy

import (
	"github.com/KubeOperator/KubeOperator/pkg/middleware"
	"github.com/KubeOperator/KubeOperator/pkg/service"
	"github.com/kataras/iris/v12"
)

var (
	keyPrefix           = "Bearer"
	AuthorizationHeader = "Authorization"
	clusterService      = service.NewClusterService()
	clusterToolService  = service.NewClusterToolService()
)

func RegisterProxy(parent iris.Party) {
	proxy := parent.Party("/proxy")
	proxy.Use(middleware.SessionMiddleware)
	proxy.Any("/kubernetes/{cluster_name}/{p:path}", KubernetesClientProxy)
	proxy.Any("/logging/{cluster_name}/{p:path}", LoggingProxy)
	proxy.Any("/loki/{cluster_name}/{p:path}", LokiProxy)
	proxy.Any("/prometheus/{cluster_name}/{p:path}", PrometheusProxy)
}
