package dto

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	clusterUtil "github.com/KubeOperator/KubeOperator/pkg/util/cluster"
)

type Cluster struct {
	model.Cluster
	NodeSize               int    `json:"nodeSize"`
	ProjectName            string `json:"projectName"`
	MultiClusterRepository string `json:"multiClusterRepository"`
}

type ClusterPage struct {
	Items []Cluster `json:"items"`
	Total int       `json:"total"`
}

type ClusterSecret struct {
	model.ClusterSecret
}

type ClusterNode struct {
	model.ClusterNode
}

type NodeCreate struct {
	HostName string `json:"hostName"`
	Role     string `json:"role"`
}

type ClusterCreate struct {
	Name                     string       `json:"name" binding:"required"`
	NodeNameRule             string       `json:"nodeNameRule" binding:"required"`
	Version                  string       `json:"version" binding:"required"`
	Provider                 string       `json:"provider"`
	Plan                     string       `json:"plan"`
	WorkerAmount             int          `json:"workerAmount"`
	NetworkType              string       `json:"networkType"`
	CiliumVersion            string       `json:"ciliumVersion"`
	CiliumTunnelMode         string       `json:"ciliumTunnelMode"`
	CiliumNativeRoutingCidr  string       `json:"ciliumNativeRoutingCidr"`
	RuntimeType              string       `json:"runtimeType"`
	DockerStorageDir         string       `json:"dockerStorageDir"`
	ContainerdStorageDir     string       `json:"containerdStorageDir"`
	FlannelBackend           string       `json:"flannelBackend"`
	CalicoIpv4PoolIpip       string       `json:"calicoIpv4PoolIpip"`
	KubeProxyMode            string       `json:"kubeProxyMode"`
	NodeportAddress          string       `json:"nodeportAddress"`
	KubeServiceNodePortRange string       `json:"kubeServiceNodePortRange"`
	EnableDnsCache           string       `json:"enableDnsCache"`
	DnsCacheVersion          string       `json:"dnsCacheVersion"`
	IngressControllerType    string       `json:"ingressControllerType"`
	Architectures            string       `json:"architectures"`
	KubeDnsDomain            string       `json:"kubeDnsDomain"`
	KubernetesAudit          string       `json:"kubernetesAudit"`
	DockerSubnet             string       `json:"dockerSubnet"`
	Nodes                    []NodeCreate `json:"nodes"`
	ProjectName              string       `json:"projectName"`
	HelmVersion              string       `json:"helmVersion"`
	NetworkInterface         string       `json:"networkInterface"`
	NetworkCidr              string       `json:"networkCidr"`
	YumOperate               string       `json:"yumOperate"`
	LbMode                   string       `json:"lbMode"`
	LbKubeApiserverIp        string       `json:"lbKubeApiserverIp"`
	KubeApiServerPort        int          `json:"kubeApiServerPort"`
	MasterScheduleType       string       `json:"masterScheduleType"`

	KubePodSubnet     string `json:"kubePodSubnet"`
	MaxNodePodNum     int    `json:"maxNodePodNum"`
	MaxNodeNum        int    `json:"maxNodeNum"`
	KubeServiceSubnet string `json:"kubeServiceSubnet"`
}

type ClusterBatch struct {
	Items     []Cluster
	Operation string
}

type Endpoint struct {
	Address string
	Port    int
}

type ClusterWithEndpoint struct {
	Cluster  model.Cluster
	Endpoint Endpoint
}

type WebkubectlToken struct {
	Token string `json:"token"`
}

type IsClusterNameExist struct {
	IsExist bool `json:"isExist"`
}

type ClusterUpgrade struct {
	ClusterName string `json:"clusterName"`
	Version     string `json:"version"`
}

type ClusterHealth struct {
	Level string              `json:"level"`
	Hooks []ClusterHealthHook `json:"hooks"`
}

type ClusterHealthHook struct {
	Name        string `json:"name"`
	Level       string `json:"level"`
	Msg         string `json:"msg"`
	AdjustValue string `json:"adjustValue"`
}

type ClusterRecoverItem struct {
	Name   string `json:"name"`
	Method string `json:"method"`
	Result string `json:"result"`
	Msg    string `json:"msg"`
}
type ClusterInfo struct {
	Name     string `json:"name"`
	Provider string `json:"provider"`
}

type ClusterLoad struct {
	Name          string `json:"name"`
	ApiServer     string `json:"apiServer"`
	Router        string `json:"router"`
	Token         string `json:"token"`
	Architectures string `json:"architectures"`
}

type ClusterLoadInfo struct {
	Name                     string `json:"name"`
	NodeNameRule             string `json:"nodeNameRule"`
	Version                  string `json:"version"`
	Architectures            string `json:"architectures"`
	LbMode                   string `json:"lbMode"`
	LbKubeApiserverIp        string `json:"lbKubeApiserverIp"`
	KubeApiServerPort        int    `json:"kubeApiServerPort"`
	KubeServiceNodePortRange string `json:"kubeServiceNodePortRange"`
	NodeportAddress          string `json:"nodeportAddress"`
	KubeProxyMode            string `json:"kubeProxyMode"`
	NetworkType              string `json:"networkType"`
	IngressControllerType    string `json:"ingressControllerType"`
	EnableDnsCache           string `json:"enableDnsCache"`
	KubeDnsDomain            string `json:"kubeDnsDomain"`
	KubernetesAudit          string `json:"kubernetesAudit"`
	RuntimeType              string `json:"runtimeType"`
	MasterScheduleType       string `json:"masterScheduleType"`

	KubePodSubnet         string         `json:"kubePodSubnet"`
	KubeServiceSubnet     string         `json:"kubeServiceSubnet"`
	MaxNodeNum            int            `json:"maxNodeNum"`
	MaxNodePodNum         int            `json:"maxNodePodNum"`
	KubeMaxPods           int            `json:"kubeMaxPods"`
	KubeNetworkNodePrefix int            `json:"kubeNetworkNodePrefix"`
	Nodes                 []NodesFromK8s `json:"nodes"`

	CephFsStatus    string                          `json:"cephFsStatus"`
	CephBlockStatus string                          `json:"cephBlockStatus"`
	NfsProvisioners []ClusterStorageProvisionerLoad `json:"nfsProvisioners"`
}

type NodesFromK8s struct {
	Name         string `json:"name"`
	Port         int    `json:"port"`
	Ip           string `json:"ip"`
	Architecture string `json:"architecture"`
	Role         string `json:"role"`
	CredentialID string `json:"credentialID"`
}

func (c ClusterCreate) ClusterCreateDto2Mo() *model.Cluster {
	cluster := model.Cluster{
		Name:          c.Name,
		NodeNameRule:  c.NodeNameRule,
		Source:        constant.ClusterSourceLocal,
		Architectures: c.Architectures,
		Provider:      c.Provider,
		Version:       c.Version,
		Status:        constant.StatusWaiting,
	}
	cluster.SpecNetwork = model.ClusterSpecNetwork{
		NetworkType:             c.NetworkType,
		CiliumVersion:           c.CiliumVersion,
		CiliumTunnelMode:        c.CiliumTunnelMode,
		CiliumNativeRoutingCidr: c.CiliumNativeRoutingCidr,
		FlannelBackend:          c.FlannelBackend,
		CalicoIpv4PoolIpip:      c.CalicoIpv4PoolIpip,
		NetworkInterface:        c.NetworkInterface,
		NetworkCidr:             c.NetworkCidr,

		Status: constant.StatusRunning,
	}
	cluster.SpecRuntime = model.ClusterSpecRuntime{
		RuntimeType:          c.RuntimeType,
		DockerStorageDir:     c.DockerStorageDir,
		ContainerdStorageDir: c.ContainerdStorageDir,
		DockerSubnet:         c.DockerSubnet,

		HelmVersion: c.HelmVersion,

		Status: constant.StatusRunning,
	}
	cluster.SpecConf = model.ClusterSpecConf{
		YumOperate: c.YumOperate,

		MaxNodeNum:        c.MaxNodeNum,
		WorkerAmount:      c.WorkerAmount,
		KubePodSubnet:     c.KubePodSubnet,
		KubeServiceSubnet: c.KubeServiceSubnet,

		KubeProxyMode:            c.KubeProxyMode,
		KubeDnsDomain:            c.KubeDnsDomain,
		KubernetesAudit:          c.KubernetesAudit,
		NodeportAddress:          c.NodeportAddress,
		KubeServiceNodePortRange: c.KubeServiceNodePortRange,

		MasterScheduleType: c.MasterScheduleType,
		LbMode:             c.LbMode,
		LbKubeApiserverIp:  c.LbKubeApiserverIp,
		KubeApiServerPort:  c.KubeApiServerPort,

		Status: constant.StatusRunning,
	}

	cluster.TaskLog = model.TaskLog{
		Type:  constant.TaskLogTypeClusterCreate,
		Phase: constant.StatusWaiting,
	}
	cluster.Secret = model.ClusterSecret{
		KubeadmToken: clusterUtil.GenerateKubeadmToken(),
	}

	nodeMask := clusterUtil.GetNodeCIDRMaskSize(c.MaxNodePodNum)
	cluster.SpecConf.KubeMaxPods = clusterUtil.MaxNodePodNumMap[nodeMask]
	cluster.SpecConf.KubeNetworkNodePrefix = nodeMask
	return &cluster
}
