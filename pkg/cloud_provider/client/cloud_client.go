package client

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
)

type CloudClient interface {
	ListDatacenter() ([]string, error)
	ListClusters() ([]interface{}, error)
	ListTemplates() ([]interface{}, error)
	ListFlavors() ([]interface{}, error)
	GetIpInUsed(network string) ([]string, error)
	UploadImage() error
	DefaultImageExist() (bool, error)
	CreateDefaultFolder() error
}

func NewCloudClient(vars map[string]interface{}) CloudClient {
	switch vars["provider"] {
	case constant.OpenStack:
		return NewOpenStackClient(vars)
	case constant.VSphere:
		return NewVSphereClient(vars)
	case constant.FusionCompute:
		return NewFusionComputeClient(vars)
	}
	return nil
}
