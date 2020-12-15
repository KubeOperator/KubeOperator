package dto

import "github.com/KubeOperator/KubeOperator/pkg/model"

type IpPool struct {
	model.IpPool
}

type IpPoolCreate struct {
	Name string `json:"name"`
}

type IpPoolOp struct {
	Operation string   `json:"operation"  validate:"required"`
	Items     []IpPool `json:"items"  validate:"required"`
}
