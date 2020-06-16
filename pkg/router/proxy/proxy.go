package proxy

import (
	"github.com/kataras/iris/v12"
)

func RegisterProxy(parent iris.Party) {
	proxy := parent.Party("/proxy")
	proxy.Any("/kubernetes/{cluster_name}/{p:path}", KubernetesClientProxy)
	proxy.Any("/logging/{cluster_name}/{p:path}", LoggingProxy)
}
