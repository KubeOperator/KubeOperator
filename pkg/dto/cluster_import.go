package dto

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
	Version                  string `json:"version" binding:"required"`
	Provider                 string `json:"provider"`
	NetworkType              string `json:"networkType"`
	CiliumVersion            string `json:"ciliumVersion"`
	CiliumTunnelMode         string `json:"ciliumTunnelMode"`
	CiliumNativeRoutingCidr  string `json:"ciliumNativeRoutingCidr"`
	RuntimeType              string `json:"runtimeType"`
	DockerStorageDIr         string `json:"dockerStorageDIr"`
	ContainerdStorageDIr     string `json:"containerdStorageDIr"`
	FlannelBackend           string `json:"flannelBackend"`
	CalicoIpv4poolIpip       string `json:"calicoIpv4PoolIpip"`
	KubeProxyMode            string `json:"kubeProxyMode"`
	NodeportAddress          string `json:"nodeportAddress"`
	KubeServiceNodePortRange string `json:"kubeServiceNodePortRange"`
	EnableDnsCache           string `json:"enableDnsCache"`
	DnsCacheVersion          string `json:"dnsCacheVersion"`
	IngressControllerType    string `json:"ingressControllerType"`
	Architectures            string `json:"architectures"`
	KubernetesAudit          string `json:"kubernetesAudit"`
	DockerSubnet             string `json:"dockerSubnet"`
	HelmVersion              string `json:"helmVersion"`
	NetworkInterface         string `json:"networkInterface"`
	NetworkCidr              string `json:"networkCidr"`
	SupportGpu               string `json:"supportGpu"`
	YumOperate               string `json:"yumOperate"`
	LbMode                   string `json:"lbMode"`
	LbKubeApiserverIp        string `json:"lbKubeApiserverIp"`
	KubeApiServerPort        int    `json:"kubeApiserverPort"`

	KubePodSubnet         string `json:"kubePodSubnet"`
	MaxNodePodNum         int    `json:"maxNodePodNum"`
	MaxNodeNum            int    `json:"maxNodeNum"`
	KubeMaxPods           int    `json:"kubeMaxPods"`
	KubeNetworkNodePrefix int    `json:"kubeNetworkNodePrefix"`
	KubeServiceSubnet     string `json:"kubeServiceSubnet"`

	Nodes []NodesFromK8s `json:"nodes"`
}
