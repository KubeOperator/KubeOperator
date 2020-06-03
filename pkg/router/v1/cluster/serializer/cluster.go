package serializer

import (
	clusterModel "github.com/KubeOperator/KubeOperator/pkg/model/cluster"
	"time"
)

type Cluster struct {
	Name     string    `json:"name"`
	Spec     Spec      `json:"spec"`
	Status   string    `json:"status"`
	NodeSize int       `json:"nodeSize"`
	CreateAt time.Time `json:"createAt"`
	UpdateAt time.Time `json:"updateAt"`
}

type Spec struct {
	Version     string `json:"version"`
	NetworkType string `json:"networkType"`
	RuntimeType string `json:"runtimeType"`
	ClusterCIDR string `json:"clusterCIDR"`
	ServiceCIDR string `json:"serviceCIDR"`
}

type Condition struct {
	Message string `json:"message"`
	Status  string `json:"status"`
	Name    string `json:"name"`
}

type Status struct {
	Phase      string      `json:"phase"`
	Conditions []Condition `json:"conditions"`
}

type Node struct {
	Role     string `json:"role"`
	Name     string `json:"name"`
	HostName string `json:"hostName"`
}

func FromModel(model clusterModel.Cluster) Cluster {
	return Cluster{
		Name: model.Name,
		Spec: Spec{
			Version: model.Spec.Version,
		},
		NodeSize: len(model.Nodes),
		Status:   model.Status.Phase,
		CreateAt: model.CreatedAt,
		UpdateAt: model.UpdatedAt,
	}
}
func ToModel(c Cluster) clusterModel.Cluster {
	return clusterModel.Cluster{
		Name: c.Name,
		Spec: clusterModel.Spec{
			Version: c.Spec.Version,
		},
	}
}

func FromNodeModel(node clusterModel.Node) Node {
	return Node{
		Role: node.Role,
		Name: node.Name,
	}
}

type ListClusterResponse struct {
	Items []Cluster `json:"items"`
	Total int       `json:"total"`
}

type ListNodeResponse struct {
	Items []Node `json:"items"`
	Total int    `json:"total"`
}

type GetClusterResponse struct {
	Item Cluster `json:"item"`
}

type CreateClusterRequest struct {
	Name        string `json:"name" binding:"required"`
	Version     string `json:"version" binding:"required"`
	NetworkType string `json:"networkType"`
	RuntimeType string `json:"runtimeType"`
	ClusterCIDR string `json:"clusterCIDR"`
	ServiceCIDR string `json:"serviceCIDR"`
	Nodes       []Node `json:"nodes"`
}

type DeleteClusterRequest struct {
	Name string `json:"name"`
}

type DeleteClusterResponse struct{}

type UpdateClusterRequest struct {
	Item Cluster `json:"item" binding:"required"`
}

type BatchClusterRequest struct {
	Operation string    `json:"operation" binding:"required"`
	Items     []Cluster `json:"items" binding:"required"`
}

type BatchClusterResponse struct {
	Items []Cluster `json:"items"`
}

type ClusterStatusResponse struct {
	Status Status `json:"status"`
}

type InitClusterResponse struct {
	Message string
}
