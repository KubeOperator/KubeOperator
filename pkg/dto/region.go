package dto

import "github.com/KubeOperator/KubeOperator/pkg/model"

type Region struct {
	model.Region
	RegionVars map[string]interface{} `json:"regionVars"`
}

type RegionCreate struct {
	Name       string                 `json:"name" validate:"koname,required"`
	Provider   string                 `json:"provider" validate:"oneof=OpenStack vSphere FusionCompute"`
	RegionVars map[string]interface{} `json:"regionVars" validate:"required"`
	Datacenter string                 `json:"datacenter" validate:"required"`
}
type RegionDatacenterRequest struct {
	RegionVars map[string]interface{} `json:"regionVars" validate:"required"`
}

type RegionOp struct {
	Operation string   `json:"operation" validate:"required"`
	Items     []Region `json:"items" validate:"required"`
}

type CloudRegionResponse struct {
	Result interface{} `json:"result"`
}
