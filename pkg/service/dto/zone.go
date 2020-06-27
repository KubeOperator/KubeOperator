package dto

import "github.com/KubeOperator/KubeOperator/pkg/model"

type Zone struct {
	model.Zone
}

type ZoneCreate struct {
	Name      string `json:"name" validate:"required"`
	Vars      string `json:"vars" validate:"required"`
	CloudZone string `json:"cloudZone" validate:"required"`
	RegionID  string `json:"regionID" validate:"required"`
}

type ZoneOp struct {
	Operation string `json:"operation" validate:"required"`
	Items     []Zone `json:"items" validate:"required"`
}

type CloudZoneResponse struct {
	Result interface{}
}
