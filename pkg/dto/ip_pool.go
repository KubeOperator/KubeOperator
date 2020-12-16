package dto

import "github.com/KubeOperator/KubeOperator/pkg/model"

type IpPool struct {
	model.IpPool
}

type IpPoolCreate struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Subnet      string `json:"subnet"`
	IpStart     string `json:"ipStart"`
	IpEnd       string `json:"ipEnd"`
	Gateway     string `json:"gateway"`
	DNS1        string `json:"dns1"`
	DNS2        string `json:"dns2"`
}

type IpPoolOp struct {
	Operation string   `json:"operation"  validate:"required"`
	Items     []IpPool `json:"items"  validate:"required"`
}
