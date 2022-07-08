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
	CalicoIpv4PoolIpip      string `json:"calicoIpv4PoolIpip"`
	NetworkInterface        string `json:"networkInterface"`
	NetworkCidr             string `json:"networkCidr"`

	Status  string `json:"status"`
	Message string `json:"message" gorm:"type:text(65535)"`
}

type ClusterSpecRuntime struct {
	common.BaseModel
	ID        string `json:"-"`
	ClusterID string `json:"-"`

	RuntimeType          string `json:"runtimeType"`
	DockerStorageDir     string `json:"dockerStorageDir"`
	ContainerdStorageDir string `json:"containerdStorageDir"`
	DockerSubnet         string `json:"dockerSubnet"`
	HelmVersion          string `json:"helmVersion"`

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
	KubernetesAudit          string `json:"kubernetesAudit"`
	NodeportAddress          string `json:"nodeportAddress"`
	KubeServiceNodePortRange string `json:"kubeServiceNodePortRange"`

	EtcdDataDir             string `json:"etcdDataDir"`
	EtcdSnapshotCount       int    `json:"etcdSnapshotCount"`
	EtcdCompactionRetention int    `json:"etcdCompactionRetention"`
	EtcdMaxRequest          int    `json:"etcdMaxRequest"`
	EtcdQuotaBackend        int    `json:"etcdQuotaBackend"`

	MasterScheduleType string `json:"masterScheduleType"`
	LbMode             string `json:"lbMode"`
	LbKubeApiserverIp  string `json:"lbKubeApiserverIp"`
	KubeApiServerPort  int    `json:"kubeApiServerPort"`
	KubeRouter         string `json:"kubeRouter"`
	AuthenticationMode string `json:"authenticationMode"`

	Status  string `json:"status"`
	Message string `json:"message" gorm:"type:text(65535)"`
}

type ClusterSpecComponent struct {
	common.BaseModel
	ID        string `json:"-"`
	ClusterID string `json:"-"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	Version   string `json:"version"`
	Describe  string `json:"describe"`
	Vars      string `json:"-"  gorm:"type:text(65535)"`

	Status  string `json:"status"`
	Message string `json:"message" gorm:"type:text(65535)"`
}

func (s *ClusterSpecConf) BeforeCreate() (err error) {
	s.ID = uuid.NewV4().String()
	return nil
}

func (s *ClusterSpecRuntime) BeforeCreate() (err error) {
	s.ID = uuid.NewV4().String()
	return nil
}

func (s *ClusterSpecNetwork) BeforeCreate() (err error) {
	s.ID = uuid.NewV4().String()
	return nil
}

func (s *ClusterSpecComponent) BeforeCreate() (err error) {
	s.ID = uuid.NewV4().String()
	return nil
}
