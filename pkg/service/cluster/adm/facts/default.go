package facts

const (
	ClusterNameFactName                  = "cluster_name"
	KubeVersionFactName                  = "kube_version"
	KubeUpgradeVersionFactName           = "kube_upgrade_version"
	ContainerRuntimeFactName             = "container_runtime"
	LbModeFactName                       = "lb_mode"
	LbKubeApiserverPortFactName          = "lb_kube_apiserver_port"
	KubeDnsDomainFactName                = "kube_dns_domain"
	KubePodSubnetFactName                = "kube_pod_subnet"
	KubeServiceSubnetFactName            = "kube_service_subnet"
	KubeNetworkNodePrefixFactName        = "kube_network_node_prefix"
	KubeMaxPodsFactName                  = "kube_max_pods"
	KubeServiceNodePortRangeFactName     = "kube_service_node_port_range"
	KubeProxyModeFactName                = "kube_proxy_mode"
	NodePortAddressFactName              = "nodeport_address"
	EnableDnsCacheFactName               = "enable_dns_cache"
	DnsCacheVersionFactName              = "dns_cache_version"
	NetworkPluginFactName                = "network_plugin"
	CiliumVersionFactName                = "cilium_version"
	CiliumTunnelModeFactName             = "cilium_tunnel_mode"
	CiliumNativeRoutingCidrFactName      = "cilium_native_routing_cidr"
	KubeImageRepositoryFactName          = "kube_image_repository"
	PodInfraContainerImageFactName       = "pod_infra_container_image"
	CertsExpiredFactName                 = "certs_expired"
	KubeCpuReservedFactName              = "kube_cpu_reserved"
	KubeMemoryReservedFactName           = "kube_memory_reserved"
	KubeEphemeralStorageReservedFactName = "kube_ephemeral_storage_reserved"
	EvictionHardImagefsAvailableFactName = "eviction_hard_imagefs_available"
	EvictionHardMemoryAvailableFactName  = "eviction_hard_memory_available"
	EvictionHardNodefsAvailableFactName  = "eviction_hard_nodefs_available"
	EvictionHardNodefsInodesFreeFactName = "eviction_hard_nodefs_inodes_free"
	KubeletRootDirFactName               = "kubelet_root_dir"
	DockerStorageDirFactName             = "docker_storage_dir"
	DockerSubnetFactName                 = "docker_subnet"
	ContainerdStorageDirFactName         = "containerd_storage_dir"
	EtcdDataDirFactName                  = "etcd_data_dir"
	BinDirFactName                       = "bin_dir"
	BaseDirFactName                      = "base_dir"
	RegistryHostnameFactName             = "registry_hostname"
	RegistryProtocolFactName             = "registry_protocol"
	RepoPortFactName                     = "repo_port"
	RegistryPortFactName                 = "registry_port"
	CorednsImageFactName                 = "coredns_image"
	KubeadmTokenFactName                 = "kubeadm_token"
	CalicoIpv4poolIpIpFactName           = "calico_ipv4pool_ipip"
	FlannelBackendFactName               = "flannel_backend"
	KubernetesAuditFactName              = "kubernetes_audit"
	IngressControllerTypeFactName        = "ingress_controller_type"
	HelmVersionFactName                  = "helm_version"
	EtcdVersionFactName                  = "etcd_version"
	DockerVersionFactName                = "docker_version"
	ContainerdVersionFactName            = "containerd_version"
	FlannelVersionFactName               = "flannel_version"
	CalicoVersionFactName                = "calico_version"
	CorednsVersionFactName               = "coredns_version"
	HelmV2VersionFactName                = "helm_v2_version"
	HelmV3VersionFactName                = "helm_v3_version"
	NginxIngressVersionFactName          = "nginx_ingress_version"
	TraefikIngressVersionFactName        = "traefik_ingress_version"
	MetricsServerVersionFactName         = "metrics_server_version"
	NetworkInterfaceFactName             = "network_interface"
	SupportGpuName                       = "support_gpu"
	YumRepoFactName                      = "yum_operate"
	NtpServerName                        = "ntp_server"
)

var DefaultFacts = map[string]string{
	KubeVersionFactName:                  "v1.18.6",
	ContainerRuntimeFactName:             "docker",
	CorednsImageFactName:                 "docker.io/kubeoperator/coredns:1.6.7",
	LbModeFactName:                       "haproxy",
	LbKubeApiserverPortFactName:          "8443",
	KubeDnsDomainFactName:                "cluster.local",
	KubePodSubnetFactName:                "10.244.0.0/18",
	KubeServiceSubnetFactName:            "10.244.64.0/18",
	DockerSubnetFactName:                 "172.17.0.1/16",
	KubeNetworkNodePrefixFactName:        "24",
	KubeMaxPodsFactName:                  "110",
	KubeProxyModeFactName:                "iptables",
	NodePortAddressFactName:              "",
	EnableDnsCacheFactName:               "disable",
	DnsCacheVersionFactName:              "1.17.0",
	NetworkPluginFactName:                "calico",
	CiliumVersionFactName:                "v1.9.5",
	CiliumTunnelModeFactName:             "vxlan",
	CiliumNativeRoutingCidrFactName:      "10.244.0.0/18",
	KubeServiceNodePortRangeFactName:     "30000-32767",
	KubeletRootDirFactName:               "/var/lib/kubelet",
	DockerStorageDirFactName:             "/var/lib/docker",
	ContainerdStorageDirFactName:         "/var/lib/containerd",
	EtcdDataDirFactName:                  "/var/lib/etcd",
	BinDirFactName:                       "/usr/local/bin",
	BaseDirFactName:                      "/opt/kubeoperator",
	RegistryHostnameFactName:             "172.16.10.64",
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
	CalicoIpv4poolIpIpFactName:           "Always",
	KubernetesAuditFactName:              "no",
	IngressControllerTypeFactName:        "nginx",
	FlannelBackendFactName:               "vxlan",
	HelmVersionFactName:                  "v3",
	EtcdVersionFactName:                  "v3.4.9",
	DockerVersionFactName:                "19.03.9",
	ContainerdVersionFactName:            "1.3.6",
	FlannelVersionFactName:               "v0.12.0",
	CalicoVersionFactName:                "v3.14.1",
	CorednsVersionFactName:               "1.6.7",
	HelmV2VersionFactName:                "v2.16.9",
	HelmV3VersionFactName:                "v3.2.4",
	NginxIngressVersionFactName:          "0.33.0",
	TraefikIngressVersionFactName:        "v2.2.1",
	MetricsServerVersionFactName:         "v0.3.6",
	NetworkInterfaceFactName:             "",
	SupportGpuName:                       "disable",
	YumRepoFactName:                      "replace",
	NtpServerName:                        "ntp1.aliyun.com",
}
