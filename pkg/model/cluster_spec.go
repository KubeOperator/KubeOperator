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
	FlannelBackend       string `json:"flannelBackend"`
	CalicoIpv4poolIpip   string `json:"calicoIpv4PoolIpip"`
	RuntimeType          string `json:"runtimeType"`
	DockerStorageDir     string `json:"dockerStorageDir"`
	ContainerdStorageDir string `json:"containerdStorageDir"`
	LbKubeApiserverIp    string `json:"lbKubeApiserverIp"`
	KubeApiServerPort    int    `json:"kubeApiServerPort"`
	KubeRouter           string `json:"kubeRouter"`
	KubePodSubnet        string `json:"kubePodSubnet"`
	KubeServiceSubnet    string `json:"kubeServiceSubnet"`
	WorkerAmount         int    `json:"workerAmount"`
}

func (s *ClusterSpec) BeforeCreate() (err error) {
	s.ID = uuid.NewV4().String()
	return nil
}

func (s ClusterSpec) TableName() string {
	return "ko_cluster_spec"
}
