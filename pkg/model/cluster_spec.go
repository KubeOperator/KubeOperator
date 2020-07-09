package model

import (
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	uuid "github.com/satori/go.uuid"
)

type ClusterSpec struct {
	common.BaseModel
	ID                   string `json:"_"`
	Version              string `json:"version"`
	Provider             string `json:"provider"`
	NetworkType          string `json:"networkType"`
	RuntimeType          string `json:"runtimeType"`
	DockerStorageDir     string `json:"dockerStorageDir"`
	ContainerdStorageDir string `json:"containerdStorageDir"`
	LbKubeApiserverIp    string `json:"lbKubeApiserverIp"`
	KubeApiServerPort    int    `json:"kubeApiServerPort"`
	KubeRouter           string `json:"kubeRouter"`
	KubePodSubnet          string `json:"clusterCidr" gorm:"column:cluster_cidr"`
	KubeServiceSubnet          string `json:"serviceCidr" gorm:"column:service_cidr"`
}

func (s *ClusterSpec) BeforeCreate() (err error) {
	s.ID = uuid.NewV4().String()
	return nil
}

func (s ClusterSpec) TableName() string {
	return "ko_cluster_spec"
}
