package serializer

import (
	clusterModel "github.com/KubeOperator/KubeOperator/pkg/model/cluster"
)

type Cluster struct {
	Name string `json:"name"`
	Spec Spec   `json:"spec"`
}

type Spec struct {
	Version string `json:"version"`
}

func FromModel(model clusterModel.Cluster) Cluster {
	return Cluster{
		Name: model.Name,
		Spec: Spec{
			Version: model.Spec.Version,
		},
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
	Name    string `json:"name" binding:"required"`
	Version string `json:"version" binding:"required"`
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
