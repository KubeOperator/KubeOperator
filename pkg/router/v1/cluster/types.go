package cluster

import "ko3-gin/pkg/model/cluster"

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

type ListResponse struct {
	items []Cluster
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
	Item Cluster
}

type UpdateRequest struct {
	Name string
}

type UpdateResponse struct {
	Item Cluster
}
