package dto

import "github.com/KubeOperator/KubeOperator/pkg/model"

type Cluster struct {
	model.Cluster
	NodeSize int    `json:"nodeSize"`
	Status   string `json:"status"`
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

type ClusterMonitor struct {
	model.ClusterMonitor
}

type ClusterSpec struct {
	model.ClusterSpec
}

type NodeCreate struct {
	HostName string `json:"hostName"`
	Role     string `json:"role"`
}

type ClusterCreate struct {
	Name                 string       `json:"name" binding:"required"`
	Version              string       `json:"version" binding:"required"`
	NetworkType          string       `json:"networkType"`
	RuntimeType          string       `json:"runtimeType"`
	DockerStorageDIr     string       `json:"dockerStorageDIr"`
	ContainerdStorageDIr string       `json:"containerdStorageDIr"`
	AppDomain            string       `json:"appDomain"`
	ClusterCIDR          string       `json:"clusterCIDR"`
	ServiceCIDR          string       `json:"serviceCIDR"`
	Nodes                []NodeCreate `json:"nodes"`
}

type ClusterBatch struct {
	Items     []Cluster
	Operation string
}

type Endpoint struct {
	Address string
	Port    int
}
