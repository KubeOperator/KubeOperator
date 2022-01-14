package dto

import "github.com/KubeOperator/KubeOperator/pkg/model"

type Ip struct {
	model.Ip
	IpPool model.IpPool `json:"ipPool"`
}

type IpCreate struct {
	IpStart    string `json:"ipStart" validate:"required,koip"`
	IpEnd      string `json:"ipEnd" validate:"required,koip"`
	Subnet     string `json:"subnet"`
	Gateway    string `json:"gateway"`
	DNS1       string `json:"dns1"`
	DNS2       string `json:"dns2"`
	IpPoolName string `json:"ipPoolName" validate:"required"`
}

type IpOp struct {
	Operation string `json:"operation"  validate:"required"`
	Items     []Ip   `json:"items"  validate:"required"`
}

type IpUpdate struct {
	Address   string `json:"address"`
	Operation string `json:"operation"`
}

type IpSync struct {
	IpPoolName string `json:"ipPoolName"`
}
