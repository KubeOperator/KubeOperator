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
	proxy.Any("/loki/{cluster_name}/{p:path}", LokiProxy)
	proxy.Any("/grafana/{cluster_name}/{p:path}", GrafanaProxy)
	proxy.Any("/prometheus/{cluster_name}/{p:path}", PrometheusProxy)
	proxy.Any("/chartmuseum/{cluster_name}/{p:path}", ChartmuseumProxy)
	proxy.Any("/dashboard/{cluster_name}/{p:path}", DashboardProxy)
	proxy.Any("/registry/{cluster_name}/{p:path}", RegistryProxy)
	proxy.Any("/kubeapps/{cluster_name}/{p:path}", KubeappsProxy)
}
