package dto

import "github.com/KubeOperator/KubeOperator/pkg/model"

type ClusterResource struct {
	model.ClusterResource
	ResourceName string `json:"resourceName"`
}

type ClusterResourceCreate struct {
	ResourceType string   `json:"resourceType" validate:"required"`
	Names        []string `json:"names" validate:"required"`
}
