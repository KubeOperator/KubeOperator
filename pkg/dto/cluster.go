package dto

import "github.com/KubeOperator/KubeOperator/pkg/model"

type Cluster struct {
	model.Cluster
	NodeSize int    `json:"nodeSize"`
	Status   string `json:"status"`
	Provider string `json:"provider"`
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
	Name                  string       `json:"name" binding:"required"`
	Version               string       `json:"version" binding:"required"`
	Provider              string       `json:"provider"`
	Plan                  string       `json:"plan"`
	WorkerAmount          int          `json:"workerAmount"`
	NetworkType           string       `json:"networkType"`
	RuntimeType           string       `json:"runtimeType"`
	DockerStorageDIr      string       `json:"dockerStorageDIr"`
	ContainerdStorageDIr  string       `json:"containerdStorageDIr"`
	FlannelBackend        string       `json:"flannelBackend"`
	CalicoIpv4poolIpip    string       `json:"calicoIpv4PoolIpip"`
	KubePodSubnet         string       `json:"kubePodSubnet"`
	KubeServiceSubnet     string       `json:"kubeServiceSubnet"`
	KubeMaxPods           int          `json:"kubeMaxPods"`
	KubeProxyMode         string       `json:"kubeProxyMode"`
	IngressControllerType string       `json:"ingressControllerType"`
	Architectures         string       `json:"architectures"`
	KubernetesAudit       bool         `json:"kubernetesAudit"`
	Nodes                 []NodeCreate `json:"nodes"`
	ProjectName           string       `json:"projectName"`
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
