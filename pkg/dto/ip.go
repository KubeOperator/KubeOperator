package dto

import "github.com/KubeOperator/KubeOperator/pkg/model"

type Ip struct {
	model.Ip
	IpPool model.IpPool `json:"ipPool"`
}

type IpCreate struct {
	StartIp  string `json:"startIp"`
	EndIp    string `json:"endIp"`
	Subnet   string `json:"subnet"`
	Gateway  string `json:"gateway"`
	DNS1     string `json:"dns1"`
	DNS2     string `json:"dns2"`
	IpPoolID string `json:"ipPoolId"`
}
