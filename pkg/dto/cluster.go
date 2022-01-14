package dto

import "github.com/KubeOperator/KubeOperator/pkg/model"

type Cluster struct {
	model.Cluster
	NodeSize               int    `json:"nodeSize"`
	Status                 string `json:"status"`
	PreStatus              string `json:"preStatus"`
	Provider               string `json:"provider"`
	Architectures          string `json:"architectures"`
	MultiClusterRepository string `json:"multiClusterRepository"`
}

type ClusterPage struct {
	Items []Cluster `json:"items"`
	Total int       `json:"total"`
}

type ClusterStatus struct {
	model.ClusterStatus
}

type ClusterSecret struct {
	model.ClusterSecret
}

type ClusterNode struct {
	model.ClusterNode
}

type ClusterSpec struct {
	model.ClusterSpec
}

type NodeCreate struct {
	HostName string `json:"hostName"`
	Role     string `json:"role"`
}

type ClusterCreate struct {
	Name                    string       `json:"name" validate:"clustername,required"`
	Version                 string       `json:"version" validate:"required"`
	Provider                string       `json:"provider" validate:"oneof=bareMetal plan"`
	Plan                    string       `json:"plan" validate:"-"`
	WorkerAmount            int          `json:"workerAmount" validate:"-"`
	NetworkType             string       `json:"networkType" validate:"oneof=flannel calico cilium"`
	CiliumVersion           string       `json:"ciliumVersion" validate:"-"`
	CiliumTunnelMode        string       `json:"ciliumTunnelMode" validate:"-"`
	CiliumNativeRoutingCidr string       `json:"ciliumNativeRoutingCidr" validate:"-"`
	RuntimeType             string       `json:"runtimeType" validate:"oneof=docker containerd"`
	DockerStorageDIr        string       `json:"dockerStorageDIr" validate:"-"`
	ContainerdStorageDIr    string       `json:"containerdStorageDIr" validate:"-"`
	FlannelBackend          string       `json:"flannelBackend" validate:"-"`
	CalicoIpv4poolIpip      string       `json:"calicoIpv4PoolIpip" validate:"-"`
	KubeProxyMode           string       `json:"kubeProxyMode" validate:"oneof=iptables ipvs"`
	NodeportAddress         string       `json:"nodeportAddress" validate:"-"`
	EnableDnsCache          string       `json:"enableDnsCache" validate:"oneof=enable disable"`
	DnsCacheVersion         string       `json:"dnsCacheVersion" validate:"-"`
	IngressControllerType   string       `json:"ingressControllerType" validate:"oneof=nginx traefik"`
	Architectures           string       `json:"architectures" validate:"oneof=arm64 amd64 all"`
	KubernetesAudit         string       `json:"kubernetesAudit" validate:"oneof=enable disable"`
	DockerSubnet            string       `json:"dockerSubnet" validate:"required"`
	Nodes                   []NodeCreate `json:"nodes" validate:"-"`
	ProjectName             string       `json:"projectName" validate:"required"`
	HelmVersion             string       `json:"helmVersion" validate:"oneof=v2 v3"`
	NetworkInterface        string       `json:"networkInterface" validate:"-"`
	SupportGpu              string       `json:"supportGpu" validate:"oneof=enable disable"`
	YumOperate              string       `json:"yumOperate" validate:"oneof=replace coexist no"`
	ClusterCIDR             string       `json:"clusterCidr" validate:"required"`
	ServiceCIDR             string       `json:"serviceCidr" validate:"required"`
	MaxPodNum               int          `json:"maxPodNum" validate:"required"`
}

type ClusterBatch struct {
	Items     []Cluster
	Operation string
}

type Endpoint struct {
	Address string
	Port    int
}

type ClusterWithEndpoint struct {
	Cluster  model.Cluster
	Endpoint Endpoint
}

type WebkubectlToken struct {
	Token string `json:"token"`
}

type ClusterLog struct {
	model.ClusterLog
}

type ClusterUpgrade struct {
	ClusterName string `json:"clusterName"`
	Version     string `json:"version"`
}

type ClusterHealth struct {
	Level string              `json:"level"`
	Hooks []ClusterHealthHook `json:"hooks"`
}

type ClusterHealthHook struct {
	Name  string `json:"name"`
	Level string `json:"level"`
	Msg   string `json:"msg"`
}

type ClusterRecoverItem struct {
	Name     string `json:"name"`
	HookName string `json:"hookName"`
	Result   string `json:"result"`
	Msg      string `json:"msg"`
}
