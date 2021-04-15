package dto

import "github.com/KubeOperator/KubeOperator/pkg/model"

type Region struct {
	model.Region
	RegionVars interface{} `json:"regionVars"`
}

type RegionCreate struct {
	Name       string      `json:"name" validate:"required"`
	Provider   string      `json:"provider" validate:"required"`
	RegionVars interface{} `json:"regionVars" validate:"required"`
	Datacenter string      `json:"datacenter" validate:"required"`
}
type RegionDatacenterRequest struct {
	RegionVars interface{} `json:"regionVars" validate:"required"`
}

type RegionOp struct {
	Operation string   `json:"operation" validate:"required"`
	Items     []Region `json:"items" validate:"required"`
}

type CloudRegionResponse struct {
	Result interface{} `json:"result"`
}

type RegionUpdate struct {
	Name       string      `json:"name" validate:"required"`
	Provider   string      `json:"provider" validate:"required"`
	RegionVars interface{} `json:"regionVars" validate:"required"`
	Datacenter string      `json:"datacenter" validate:"required"`
}
