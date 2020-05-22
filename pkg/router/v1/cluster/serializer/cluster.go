package serializer

import (
	clusterModel "github.com/KubeOperator/KubeOperator/pkg/model/cluster"
	"time"
)

type Cluster struct {
	Name     string    `json:"name"`
	Spec     Spec      `json:"spec"`
	Status   Status    `json:"status"`
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

type Status struct {
	Phase string `json:"phase"`
}

func FromModel(model clusterModel.Cluster) Cluster {
	return Cluster{
		Name: model.Name,
		Spec: Spec{
			Version: model.Spec.Version,
		},
		Status: Status{
			Phase: model.Status.Phase,
		},
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

type ListClusterResponse struct {
	Items []Cluster `json:"items"`
	Total int       `json:"total"`
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
