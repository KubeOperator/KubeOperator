package cluster

import (
	"ko3-gin/pkg/model/cluster"
	"ko3-gin/pkg/model/common"
)

type Cluster struct {
	Name   string
	Status string
}

func FromModel(model cluster.Cluster) Cluster {
	return Cluster{
		Name:   model.Name,
		Status: model.Status,
	}
}

func ToModel(c Cluster) cluster.Cluster {
	return cluster.Cluster{
		BaseModel: common.BaseModel{
			Name: c.Name,
		},
	}
}

type ListResponse struct {
	Items []Cluster
}

type GetResponse struct {
	Item Cluster
}

type CreateRequest struct {
	Name string
}

type CreateResponse struct {
	Item Cluster
}

type DeleteRequest struct {
	Name string
}

type DeleteResponse struct {
}

type UpdateRequest struct {
	Name string
}

type UpdateResponse struct {
	Item Cluster
}
type BatchRequest struct {
	Operation string
	Items     []Cluster
}

type BatchResponse struct {
	Items []Cluster
}
