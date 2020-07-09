package dto

import "github.com/KubeOperator/KubeOperator/pkg/model"

type Zone struct {
	model.Zone
	CloudVars interface{} `json:"cloudVars"`
	Region    Region      `json:"region"`
}

type ZoneCreate struct {
	Name      string      `json:"name" validate:"required"`
	CloudVars interface{} `json:"cloudVars" validate:"required"`
	RegionID  string      `json:"regionID" validate:"required"`
}

type ZoneOp struct {
	Operation string `json:"operation" validate:"required"`
	Items     []Zone `json:"items" validate:"required"`
}

type CloudZoneResponse struct {
	Result interface{} `json:"result"`
}

type CloudZoneRequest struct {
	CloudVars  interface{} `json:"cloudVars" validate:"required"`
	Datacenter string      `json:"datacenter" validate:"required"`
}
