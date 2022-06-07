package model

import (
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	uuid "github.com/satori/go.uuid"
)

type ClusterSpecNetwork struct {
	common.BaseModel
	ID        string `json:"-"`
	ClusterID string `json:"-"`

	NetworkType             string `json:"networkType"`
	CiliumVersion           string `json:"ciliumVersion"`
	CiliumTunnelMode        string `json:"ciliumTunnelMode"`
	CiliumNativeRoutingCidr string `json:"ciliumNativeRoutingCidr"`
	FlannelBackend          string `json:"flannelBackend"`
	CalicoIpv4poolIpip      string `json:"calicoIpv4PoolIpip"`
	NetworkInterface        string `json:"networkInterface"`
	NetworkCidr             string `json:"networkCidr"`

	Status  string `json:"status"`
	Message string `json:"message" gorm:"type:text(65535)"`
}

type ClusterSpecRelyOn struct {
	common.BaseModel
	ID        string `json:"-"`
	ClusterID string `json:"-"`

	RuntimeType           string `json:"runtimeType"`
	DockerStorageDir      string `json:"dockerStorageDir"`
	ContainerdStorageDir  string `json:"containerdStorageDir"`
	DockerSubnet          string `json:"dockerSubnet"`
	HelmVersion           string `json:"helmVersion"`
	IngressControllerType string `json:"ingressControllerType"`

	Status  string `json:"status"`
	Message string `json:"message" gorm:"type:text(65535)"`
}

type ClusterSpecConf struct {
	common.BaseModel
	ID        string `json:"-"`
	ClusterID string `json:"-"`

	YumOperate string `json:"yumOperate"`

	MaxNodeNum            int    `json:"maxNodeNum"`
	WorkerAmount          int    `json:"workerAmount"`
	KubeMaxPods           int    `json:"kubeMaxPods"`
	KubeNetworkNodePrefix int    `json:"kubeNetworkNodePrefix"`
	KubePodSubnet         string `json:"kubePodSubnet"`
	KubeServiceSubnet     string `json:"kubeServiceSubnet"`

	KubeProxyMode            string `json:"kubeProxyMode"`
	KubeDnsDomain            string `json:"kubeDnsDomain"`
	EnableDnsCache           string `json:"enableDnsCache"`
	DnsCacheVersion          string `json:"dnsCacheVersion"`
	KubernetesAudit          string `json:"kubernetesAudit"`
	NodeportAddress          string `json:"nodeportAddress"`
	KubeServiceNodePortRange string `json:"kubeServiceNodePortRange"`

	MasterScheduleType string `json:"masterScheduleType"`
	LbMode             string `json:"lbMode"`
	LbKubeApiserverIp  string `json:"lbKubeApiserverIp"`
	KubeApiServerPort  int    `json:"kubeApiServerPort"`
	KubeRouter         string `json:"kubeRouter"`

	SupportGpu string `json:"supportGpu"`

	Status  string `json:"status"`
	Message string `json:"message" gorm:"type:text(65535)"`
}

func (s *ClusterSpecConf) BeforeCreate() (err error) {
	s.ID = uuid.NewV4().String()
	return nil
}

func (s *ClusterSpecRelyOn) BeforeCreate() (err error) {
	s.ID = uuid.NewV4().String()
	return nil
}

func (s *ClusterSpecNetwork) BeforeCreate() (err error) {
	s.ID = uuid.NewV4().String()
	return nil
}
