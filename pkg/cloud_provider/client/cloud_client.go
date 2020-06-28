package client

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
)

type CloudClient interface {
	ListZones() string
	ListDatacenter() ([]string, error)
	ListClusters(datacenter string) ([]string, error)
}

func NewCloudClient(vars map[string]interface{}) CloudClient {
	if vars["provider"] == constant.OpenStack {
		return NewOpenStackClient(vars)
	}
	if vars["provider"] == constant.VSphere {
		return NewVSphereClient(vars)
	}
	return nil
}
