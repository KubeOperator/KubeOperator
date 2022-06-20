package dto

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/model"
)

type ClusterImport struct {
	Name          string      `json:"name"`
	ApiServer     string      `json:"apiServer"`
	Router        string      `json:"router"`
	Token         string      `json:"token"`
	ProjectName   string      `json:"projectName"`
	Architectures string      `json:"architectures"`
	IsKoCluster   bool        `json:"isKoCluster"`
	KoClusterInfo clusterInfo `json:"clusterInfo"`
}

type clusterInfo struct {
	Name                     string `json:"name"`
	NodeNameRule             string `json:"nodeNameRule"`
	Version                  string `json:"version" binding:"required"`
	Provider                 string `json:"provider"`
	NetworkType              string `json:"networkType"`
	CiliumVersion            string `json:"ciliumVersion"`
	CiliumTunnelMode         string `json:"ciliumTunnelMode"`
	CiliumNativeRoutingCidr  string `json:"ciliumNativeRoutingCidr"`
	RuntimeType              string `json:"runtimeType"`
	DockerStorageDir         string `json:"dockerStorageDir"`
	ContainerdStorageDir     string `json:"containerdStorageDir"`
	FlannelBackend           string `json:"flannelBackend"`
	CalicoIpv4PoolIpip       string `json:"calicoIpv4PoolIpip"`
	KubeProxyMode            string `json:"kubeProxyMode"`
	NodeportAddress          string `json:"nodeportAddress"`
	KubeServiceNodePortRange string `json:"kubeServiceNodePortRange"`
	EnableDnsCache           string `json:"enableDnsCache"`
	DnsCacheVersion          string `json:"dnsCacheVersion"`
	IngressControllerType    string `json:"ingressControllerType"`
	Architectures            string `json:"architectures"`
	KubeDnsDomain            string `json:"kubeDnsDomain"`
	KubernetesAudit          string `json:"kubernetesAudit"`
	DockerSubnet             string `json:"dockerSubnet"`
	HelmVersion              string `json:"helmVersion"`
	NetworkInterface         string `json:"networkInterface"`
	NetworkCidr              string `json:"networkCidr"`
	YumOperate               string `json:"yumOperate"`
	LbMode                   string `json:"lbMode"`
	LbKubeApiserverIp        string `json:"lbKubeApiserverIp"`
	KubeApiServerPort        int    `json:"kubeApiServerPort"`

	KubePodSubnet         string `json:"kubePodSubnet"`
	MaxNodePodNum         int    `json:"maxNodePodNum"`
	MaxNodeNum            int    `json:"maxNodeNum"`
	KubeMaxPods           int    `json:"kubeMaxPods"`
	KubeNetworkNodePrefix int    `json:"kubeNetworkNodePrefix"`
	KubeServiceSubnet     string `json:"kubeServiceSubnet"`

	Nodes        []NodesFromK8s                  `json:"nodes"`
	Provisioners []ClusterStorageProvisionerLoad `json:"provisioners"`
}

func (c ClusterImport) ClusterImportDto2Mo() (*model.Cluster, error) {
	var (
		address string
		port    int
		cluster model.Cluster
	)
	if strings.HasSuffix(c.ApiServer, "/") {
		c.ApiServer = strings.Replace(c.ApiServer, "/", "", -1)
	}
	c.ApiServer = strings.Replace(c.ApiServer, "http://", "", -1)
	c.ApiServer = strings.Replace(c.ApiServer, "https://", "", -1)
	if !strings.Contains(c.ApiServer, ":") {
		return &cluster, fmt.Errorf("check whether apiserver(%s) has no ports", c.ApiServer)
	}
	strs := strings.Split(c.ApiServer, ":")
	address = strs[0]
	port, _ = strconv.Atoi(strs[1])

	cluster = model.Cluster{
		Name:          c.Name,
		NodeNameRule:  c.KoClusterInfo.NodeNameRule,
		Source:        constant.ClusterSourceLocal,
		Architectures: c.Architectures,
		Provider:      constant.ClusterProviderBareMetal,
		Version:       c.KoClusterInfo.Version,
	}
	cluster.TaskLog = model.TaskLog{
		Type:  constant.TaskLogTypeClusterImport,
		Phase: constant.StatusWaiting,
	}
	cluster.SpecConf = model.ClusterSpecConf{
		LbKubeApiserverIp: address,
		KubeApiServerPort: port,
		KubeRouter:        c.Router,
	}
	cluster.Secret = model.ClusterSecret{
		KubeadmToken:    "",
		KubernetesToken: c.Token,
	}

	if !c.IsKoCluster {
		return &cluster, nil
	}

	cluster.Name = c.KoClusterInfo.Name
	cluster.Source = constant.ClusterSourceKoExternal
	cluster.SpecNetwork = model.ClusterSpecNetwork{
		NetworkType:             c.KoClusterInfo.NetworkType,
		CiliumVersion:           c.KoClusterInfo.CiliumVersion,
		CiliumTunnelMode:        c.KoClusterInfo.CiliumTunnelMode,
		CiliumNativeRoutingCidr: c.KoClusterInfo.CiliumNativeRoutingCidr,
		FlannelBackend:          c.KoClusterInfo.FlannelBackend,
		CalicoIpv4PoolIpip:      c.KoClusterInfo.CalicoIpv4PoolIpip,
		NetworkInterface:        c.KoClusterInfo.NetworkInterface,
		NetworkCidr:             c.KoClusterInfo.NetworkCidr,

		Status: constant.StatusRunning,
	}
	cluster.SpecRuntime = model.ClusterSpecRuntime{
		RuntimeType:          c.KoClusterInfo.RuntimeType,
		DockerStorageDir:     c.KoClusterInfo.DockerStorageDir,
		ContainerdStorageDir: c.KoClusterInfo.ContainerdStorageDir,
		DockerSubnet:         c.KoClusterInfo.DockerSubnet,

		HelmVersion: c.KoClusterInfo.HelmVersion,

		Status: constant.StatusRunning,
	}
	cluster.SpecConf = model.ClusterSpecConf{
		YumOperate: c.KoClusterInfo.YumOperate,

		MaxNodeNum:        c.KoClusterInfo.MaxNodeNum,
		KubePodSubnet:     c.KoClusterInfo.KubePodSubnet,
		KubeServiceSubnet: c.KoClusterInfo.KubeServiceSubnet,

		KubeProxyMode:            c.KoClusterInfo.KubeProxyMode,
		KubeDnsDomain:            c.KoClusterInfo.KubeDnsDomain,
		KubernetesAudit:          c.KoClusterInfo.KubernetesAudit,
		NodeportAddress:          c.KoClusterInfo.NodeportAddress,
		KubeServiceNodePortRange: c.KoClusterInfo.KubeServiceNodePortRange,

		LbMode:            c.KoClusterInfo.LbMode,
		LbKubeApiserverIp: c.KoClusterInfo.LbKubeApiserverIp,
		KubeApiServerPort: c.KoClusterInfo.KubeApiServerPort,

		Status: constant.StatusRunning,
	}
	return &cluster, nil
}
