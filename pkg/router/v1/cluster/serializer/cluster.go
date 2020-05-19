package serializer

import (
	clusterModel "ko3-gin/pkg/model/cluster"
	"ko3-gin/pkg/model/common"
)

type Cluster struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

func FromModel(model clusterModel.Cluster) Cluster {
	return Cluster{
		Name:   model.Name,
		Status: model.Status,
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
	Name string `json:"name"`
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
	Name string `json:"name"`
}

type UpdateResponse struct {
	Item Cluster `json:"item"`
}
type BatchRequest struct {
	Operation string    `json:"operation"`
	Items     []Cluster `json:"items"`
}

type BatchResponse struct {
	Items []Cluster `json:"items"`
}
