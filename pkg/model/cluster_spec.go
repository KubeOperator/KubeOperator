package model

import (
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	uuid "github.com/satori/go.uuid"
)

type ClusterSpec struct {
	common.BaseModel
	ID                    string `json:"-"`
	Version               string `json:"version"`
	UpgradeVersion        string `json:"upgradeVersion"`
	Provider              string `json:"provider"`
	NetworkType           string `json:"networkType"`
	FlannelBackend        string `json:"flannelBackend"`
	CalicoIpv4poolIpip    string `json:"calicoIpv4PoolIpip"`
	RuntimeType           string `json:"runtimeType"`
	DockerStorageDir      string `json:"dockerStorageDir"`
	ContainerdStorageDir  string `json:"containerdStorageDir"`
	LbKubeApiserverIp     string `json:"lbKubeApiserverIp"`
	KubeApiServerPort     int    `json:"kubeApiServerPort"`
	KubeRouter            string `json:"kubeRouter"`
	KubePodSubnet         string `json:"kubePodSubnet"`
	KubeServiceSubnet     string `json:"kubeServiceSubnet"`
	DockerSubnet          string `json:"docker_subnet"`
	WorkerAmount          int    `json:"workerAmount"`
	KubeMaxPods           int    `json:"kubeMaxPods"`
	KubeProxyMode         string `json:"kubeProxyMode"`
	IngressControllerType string `json:"ingressControllerType"`
	Architectures         string `json:"architectures"`
	KubernetesAudit       string `json:"kubernetesAudit"`
	HelmVersion           string `json:"helmVersion"`
	NetworkInterface      string `json:"networkInterface"`
	SupportGpu            string `json:"supportGpu"`
	YumOperate            string `json:"yumOperate"`
}

func (s *ClusterSpec) BeforeCreate() (err error) {
	s.ID = uuid.NewV4().String()
	return nil
}
