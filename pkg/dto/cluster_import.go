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
	Router        string      `json:"router"`
	ProjectName   string      `json:"projectName"`
	Architectures string      `json:"architectures"`
	IsKoCluster   bool        `json:"isKoCluster"`
	KoClusterInfo clusterInfo `json:"clusterInfo"`

	AuthenticationMode string `json:"authenticationMode"`
	ApiServer          string `json:"apiServer"`
	Token              string `json:"token"`
	CertDataStr        string `json:"certDataStr"`
	KeyDataStr         string `json:"keyDataStr"`
	ConfigContent      string `json:"configContent"`
}

type clusterInfo struct {
	Name          string `json:"name"`
	NodeNameRule  string `json:"nodeNameRule"`
	Version       string `json:"version" binding:"required"`
	Architectures string `json:"architectures"`
	Provider      string `json:"provider"`
	YumOperate    string `json:"yumOperate"`

	NetworkType             string `json:"networkType"`
	CiliumVersion           string `json:"ciliumVersion"`
	CiliumTunnelMode        string `json:"ciliumTunnelMode"`
	CiliumNativeRoutingCidr string `json:"ciliumNativeRoutingCidr"`
	FlannelBackend          string `json:"flannelBackend"`
	CalicoIpv4PoolIpip      string `json:"calicoIpv4PoolIpip"`

	RuntimeType          string `json:"runtimeType"`
	DockerSubnet         string `json:"dockerSubnet"`
	DockerStorageDir     string `json:"dockerStorageDir"`
	ContainerdStorageDir string `json:"containerdStorageDir"`
	NetworkInterface     string `json:"networkInterface"`
	NetworkCidr          string `json:"networkCidr"`

	KubePodSubnet            string `json:"kubePodSubnet"`
	MaxNodePodNum            int    `json:"maxNodePodNum"`
	MaxNodeNum               int    `json:"maxNodeNum"`
	KubeMaxPods              int    `json:"kubeMaxPods"`
	KubeNetworkNodePrefix    int    `json:"kubeNetworkNodePrefix"`
	KubeServiceSubnet        string `json:"kubeServiceSubnet"`
	KubeProxyMode            string `json:"kubeProxyMode"`
	KubeDnsDomain            string `json:"kubeDnsDomain"`
	KubernetesAudit          string `json:"kubernetesAudit"`
	NodeportAddress          string `json:"nodeportAddress"`
	KubeServiceNodePortRange string `json:"kubeServiceNodePortRange"`

	HelmVersion             string `json:"helmVersion"`
	EtcdDataDir             string `json:"etcdDataDir"`
	EtcdSnapshotCount       int    `json:"etcdSnapshotCount"`
	EtcdCompactionRetention int    `json:"etcdCompactionRetention"`
	EtcdMaxRequest          int    `json:"etcdMaxRequest"`
	EtcdQuotaBackend        int    `json:"etcdQuotaBackend"`

	LbMode            string         `json:"lbMode"`
	LbKubeApiserverIp string         `json:"lbKubeApiserverIp"`
	KubeApiServerPort int            `json:"kubeApiServerPort"`
	Nodes             []NodesFromK8s `json:"nodes"`
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
		Status:        constant.StatusRunning,
	}
	cluster.TaskLog = model.TaskLog{
		Type:  constant.TaskLogTypeClusterImport,
		Phase: constant.StatusWaiting,
	}
	cluster.SpecConf = model.ClusterSpecConf{
		LbKubeApiserverIp:  address,
		KubeApiServerPort:  port,
		KubeRouter:         c.Router,
		AuthenticationMode: c.AuthenticationMode,
	}
	cluster.Secret = model.ClusterSecret{KubeadmToken: "",
		KubernetesToken: c.Token,
	}
	switch c.AuthenticationMode {
	case constant.AuthenticationModeBearer:
		cluster.Secret.KubernetesToken = c.Token
	case constant.AuthenticationModeCertificate:
		cluster.Secret.CertDataStr = c.CertDataStr
		cluster.Secret.KeyDataStr = c.KeyDataStr
	case constant.AuthenticationModeConfigFile:
		cluster.Secret.ConfigContent = c.ConfigContent
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

		EtcdDataDir:             c.KoClusterInfo.EtcdDataDir,
		EtcdSnapshotCount:       c.KoClusterInfo.EtcdSnapshotCount,
		EtcdCompactionRetention: c.KoClusterInfo.EtcdCompactionRetention,
		EtcdMaxRequest:          c.KoClusterInfo.EtcdMaxRequest,
		EtcdQuotaBackend:        c.KoClusterInfo.EtcdQuotaBackend,
		LbMode:                  c.KoClusterInfo.LbMode,
		LbKubeApiserverIp:       c.KoClusterInfo.LbKubeApiserverIp,
		KubeApiServerPort:       c.KoClusterInfo.KubeApiServerPort,
		AuthenticationMode:      c.AuthenticationMode,

		Status: constant.StatusRunning,
	}
	return &cluster, nil
}
