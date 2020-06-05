package facts

const (
	KubeVersionFactName              = "kube_version"
	ContainerRuntimeFactName         = "container_runtime"
	LbModeFactName                   = "lb_mode"
	LbKubeApiserverPortFactName      = "lb_kube_apiserver_port"
	KubeDnsDomainFactName            = "kube_dns_domain"
	KubePodSubnetFactName            = "kube_pod_subnet"
	KubeServiceSubnetFactName        = "kube_service_subnet"
	KubeNetworkNodePrefixFactName    = "kube_network_node_prefix"
	KubeMaxPodsFactName              = "kube_max_pods"
	KubeServiceNodePortRangeFactName = "kube_service_node_port_range"
	KubeProxyModeFactName            = "kube_proxy_mode"
	NetworkPluginFactName            = "network_plugin"
	KubeImageRepositoryFactName      = "kube_image_repository"
	PodInfraContainerImageFactName   = "pod_infra_container_image"
	CertsExpiredFactName             = "certs_expired"

	KubeCpuReservedFactName              = "kube_cpu_reserved"
	KubeMemoryReservedFactName           = "kube_memory_reserved"
	KubeEphemeralStorageReservedFactName = "kube_ephemeral_storage_reserved"

	EvictionHardImagefsAvailableFactName = "eviction_hard_imagefs_available"
	EvictionHardMemoryAvailableFactName  = "eviction_hard_memory_available"
	EvictionHardNodefsAvailableFactName  = "eviction_hard_nodefs_available"
	EvictionHardNodefsInodesFreeFactName = "eviction_hard_nodefs_inodes_free"

	KubeletRootDirFactName       = "kubelet_root_dir"
	DockerStorageDirFactName     = "docker_storage_dir"
	ContainerdStorageDirFactName = "containerd_storage_dir"

	EtcdDataDirFactName = "etcd_data_dir"

	BinDirFactName  = "bin_dir"
	BaseDirFactName = "base_dir"

	LocalHostnameFactName = "local_hostname"
	RepoPortFactName      = "repo_port"
	RegistryPortFactName  = "registry_port"
	CorednsImageFactName  = "coredns_image"
	KubeadmTokenFactName  = "kubeadm_token"
)

var DefaultFacts = map[string]string{
	KubeVersionFactName:                  "v1.18.3",
	ContainerRuntimeFactName:             "docker",
	CorednsImageFactName:                 "docker.io/kubeoperator/coredns:1.6.7",
	LbModeFactName:                       "haproxy",
	LbKubeApiserverPortFactName:          "8443",
	KubeDnsDomainFactName:                "cluster.local",
	KubePodSubnetFactName:                "10.244.0.0/18",
	KubeServiceSubnetFactName:            "10.244.64.0/18",
	KubeNetworkNodePrefixFactName:        "24",
	KubeMaxPodsFactName:                  "110",
	KubeProxyModeFactName:                "iptables",
	NetworkPluginFactName:                "calico",
	KubeServiceNodePortRangeFactName:     "30000-32767",
	KubeletRootDirFactName:               "/var/lib/kubelet",
	DockerStorageDirFactName:             "/var/lib/docker",
	ContainerdStorageDirFactName:         "/var/lib/containerd",
	EtcdDataDirFactName:                  "/var/lib/etcd",
	BinDirFactName:                       "/usr/local/bin",
	BaseDirFactName:                      "/opt/kubeoperator",
	LocalHostnameFactName:                "172.16.10.64",
	RepoPortFactName:                     "8081",
	RegistryPortFactName:                 "8082",
	KubeImageRepositoryFactName:          "docker.io/kubeoperator",
	PodInfraContainerImageFactName:       "docker.io/kubeoperator/pause:3.1",
	CertsExpiredFactName:                 "36500",
	EvictionHardImagefsAvailableFactName: "15%",
	EvictionHardMemoryAvailableFactName:  "100Mi",
	EvictionHardNodefsAvailableFactName:  "10%",
	EvictionHardNodefsInodesFreeFactName: "5%",
	KubeadmTokenFactName:                 "abcdef.0123456789abcdef",
	KubeCpuReservedFactName:              "100m",
	KubeMemoryReservedFactName:           "256M",
	KubeEphemeralStorageReservedFactName: "1G",
}
