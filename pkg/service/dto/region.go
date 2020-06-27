package dto

import "github.com/KubeOperator/KubeOperator/pkg/model"

type Region struct {
	model.Region
}

type RegionCreate struct {
	Name       string      `json:"name" validate:"required"`
	RegionVars interface{} `json:"regionVars" validate:"required"`
	Datacenter string      `json:"datacenter" validate:"required"`
	Vars       string      `json:"vars" validate:"required"`
}

type RegionOp struct {
	Operation string   `json:"operation" validate:"required"`
	Items     []Region `json:"items" validate:"required"`
}

type CloudRegionResponse struct {
	Result interface{}
}
