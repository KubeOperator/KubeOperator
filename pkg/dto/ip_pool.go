package dto

import "github.com/KubeOperator/KubeOperator/pkg/model"

type IpPool struct {
	model.IpPool
	IpUsed int `json:"ipUsed"`
}

type IpPoolCreate struct {
	Name        string `json:"name" validate:"koname,required,max=30"`
	Description string `json:"description" validate:"max=30"`
	Subnet      string `json:"subnet" validate:"required"`
	IpStart     string `json:"ipStart" validate:"required,koip"`
	IpEnd       string `json:"ipEnd" validate:"required,koip"`
	Gateway     string `json:"gateway" validate:"required"`
	DNS1        string `json:"dns1" validate:"required"`
	DNS2        string `json:"dns2" validate:"required"`
}

type IpPoolOp struct {
	Operation string   `json:"operation"  validate:"required"`
	Items     []IpPool `json:"items"  validate:"required"`
}
