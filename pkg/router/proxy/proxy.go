package proxy

import (
	"github.com/KubeOperator/KubeOperator/pkg/service"
	"github.com/kataras/iris/v12"
)

var (
	keyPrefix           = "Bearer"
	AuthorizationHeader = "Authorization"
	clusterService      = service.NewClusterService()
)

func RegisterProxy(parent iris.Party) {
	proxy := parent.Party("/proxy")
	proxy.Any("/kubernetes/{cluster_name}/{p:path}", KubernetesClientProxy)
	proxy.Any("/logging/{cluster_name}/{p:path}", LoggingProxy)
	proxy.Any("/prometheus/{p:path}", PrometheusProxy)
}
