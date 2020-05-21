package serializer

import (
	clusterModel "github.com/KubeOperator/KubeOperator/pkg/model/cluster"
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
)

type Cluster struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

func FromModel(model clusterModel.Cluster) Cluster {
	return Cluster{
		Name:   model.Name,
	}
}

func ToModel(c Cluster) clusterModel.Cluster {
	return clusterModel.Cluster{
		BaseModel: common.BaseModel{
			Name: c.Name,
		},
	}
}

type ListResponse struct {
	Items []Cluster `json:"items"`
	Total int       `json:"total"`
}

type GetResponse struct {
	Item Cluster `json:"item"`
}

type CreateRequest struct {
	Name string ` json:"name" binding:"required"`
}

type CreateResponse struct {
	Item Cluster `json:"item"`
}

type DeleteRequest struct {
	Name string `json:"name"`
}

type DeleteResponse struct {
}

type UpdateRequest struct {
	Item Cluster `json:"item" binding:"required"`
}

type UpdateResponse struct {
	Item Cluster `json:"item"`
}
type BatchRequest struct {
	Operation string    `json:"operation" binding:"required"`
	Items     []Cluster `json:"items" binding:"required"`
}

type BatchResponse struct {
	Items []Cluster `json:"items"`
}
