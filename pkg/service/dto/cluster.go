package dto

import "github.com/KubeOperator/KubeOperator/pkg/model"

type Cluster struct {
	model.Cluster
	Spec   model.ClusterSpec   `json:"spec"`
	Status model.ClusterStatus `json:"status"`
	Secret model.ClusterSecret `json:"secret"`
	Nodes  []model.ClusterNode `json:"nodes"`
}

type ClusterStatus struct {
	model.ClusterStatus
	Conditions []model.ClusterStatusCondition
}

type ClusterSecret struct {
	model.ClusterSecret
}

type ClusterCreate struct {
	Name                 string                  `json:"name" binding:"required"`
	Version              string                  `json:"version" binding:"required"`
	NetworkType          string                  `json:"networkType"`
	RuntimeType          string                  `json:"runtimeType"`
	DockerStorageDIr     string                  `json:"dockerStorageDIr"`
	ContainerdStorageDIr string                  `json:"containerdStorageDIr"`
	AppDomain            string                  `json:"appDomain"`
	ClusterCIDR          string                  `json:"clusterCIDR"`
	ServiceCIDR          string                  `json:"serviceCIDR"`
	Nodes                []struct{ Host string } `json:"nodes"`
}
