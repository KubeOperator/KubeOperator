package tools

import (
	clusterModel "github.com/KubeOperator/KubeOperator/pkg/model/cluster"
	toolModel "github.com/KubeOperator/KubeOperator/pkg/model/cluster/tool"
	"github.com/KubeOperator/KubeOperator/pkg/util/helm"
	"k8s.io/client-go/kubernetes"
	"reflect"
	"runtime"
	"strings"
)

type Handler func(*Tool) error

type CLusterToolsAdm struct {
	ingressHandlers []Handler
	efkHandlers     []Handler
}

type Tool struct {
	cluster clusterModel.Cluster
	toolModel.Tool
	Condition        string
	Values           map[string]interface{}
	HelmClient       helm.Interface
	KubernetesClient *kubernetes.Clientset
}

func NewClusterToolsAdm() (*CLusterToolsAdm, error) {
	cta := new(CLusterToolsAdm)
	cta.ingressHandlers = []Handler{
		cta.EnsureClusterReachable,
		cta.EnsureIngressInstall,
		cta.EnsureIngressRunning,
	}
	return cta, nil
}

func (h Handler) name() string {
	name := runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name()
	i := strings.Index(name, "Ensure")
	if i == -1 {
		return ""
	}
	return strings.TrimSuffix(name[i:], "-fm")
}
