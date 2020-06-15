package model

import (
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	uuid "github.com/satori/go.uuid"
)

type ClusterSpec struct {
	common.BaseModel
	ID                    string
	Version               string
	Provider              string
	NetworkType           string
	RuntimeType           string
	DockerStorageDir      string
	ContainerdStorageDir  string
	LbKubeApiserverIp     string
	AppDomain             string
	ClusterCIDR           string `gorm:"column:cluster_cidr"`
	ServiceCIDR           string `gorm:"column:service_cidr"`
}

func (s *ClusterSpec) BeforeCreate() (err error) {
	s.ID = uuid.NewV4().String()
	return nil
}

func (s ClusterSpec) TableName() string {
	return "ko_cluster_spec"
}
