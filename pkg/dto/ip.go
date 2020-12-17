package dto

import "github.com/KubeOperator/KubeOperator/pkg/model"

type Ip struct {
	model.Ip
	IpPool model.IpPool `json:"ipPool"`
}

type IpCreate struct {
	IpStart    string `json:"ipStart"`
	IpEnd      string `json:"ipEnd"`
	Subnet     string `json:"subnet"`
	Gateway    string `json:"gateway"`
	DNS1       string `json:"dns1"`
	DNS2       string `json:"dns2"`
	IpPoolName string `json:"ipPoolName"`
}

type IpOp struct {
	Operation string `json:"operation"  validate:"required"`
	Items     []Ip   `json:"items"  validate:"required"`
}
