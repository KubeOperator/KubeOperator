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
	Name                    string       `json:"name" binding:"required"`
	Version                 string       `json:"version" binding:"required"`
	Provider                string       `json:"provider"`
	Plan                    string       `json:"plan"`
	WorkerAmount            int          `json:"workerAmount"`
	NetworkType             string       `json:"networkType"`
	CiliumVersion           string       `json:"ciliumVersion"`
	CiliumTunnelMode        string       `json:"ciliumTunnelMode"`
	CiliumNativeRoutingCidr string       `json:"ciliumNativeRoutingCidr"`
	RuntimeType             string       `json:"runtimeType"`
	DockerStorageDIr        string       `json:"dockerStorageDIr"`
	ContainerdStorageDIr    string       `json:"containerdStorageDIr"`
	FlannelBackend          string       `json:"flannelBackend"`
	CalicoIpv4poolIpip      string       `json:"calicoIpv4PoolIpip"`
	KubeProxyMode           string       `json:"kubeProxyMode"`
	EnableDnsCache          string       `json:"enableDnsCache"`
	DnsCacheVersion         string       `json:"dnsCacheVersion"`
	IngressControllerType   string       `json:"ingressControllerType"`
	Architectures           string       `json:"architectures"`
	KubernetesAudit         string       `json:"kubernetesAudit"`
	DockerSubnet            string       `json:"dockerSubnet"`
	Nodes                   []NodeCreate `json:"nodes"`
	ProjectName             string       `json:"projectName"`
	HelmVersion             string       `json:"helmVersion"`
	NetworkInterface        string       `json:"networkInterface"`
	SupportGpu              string       `json:"supportGpu"`
	YumOperate              string       `json:"yumOperate"`
	ClusterCIDR             string       `json:"clusterCidr"`
	MaxClusterServiceNum    int          `json:"maxClusterServiceNum"`
	MaxNodePodNum           int          `json:"maxNodePodNum"`
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

type IsClusterNameExist struct {
	IsExist bool `json:"isExist"`
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
