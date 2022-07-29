package facts

const (
	ClusterNameFactName        = "cluster_name"
	NodeNameRuleFactName       = "node_name_rule"
	ComponentOptionFactName    = "component_created_by"
	KubeVersionFactName        = "kube_version"
	KubeUpgradeVersionFactName = "kube_upgrade_version"
	YumRepoFactName            = "yum_operate"

	KubeNetworkNodePrefixFactName    = "kube_network_node_prefix"
	KubeMaxPodsFactName              = "kube_max_pods"
	KubePodSubnetFactName            = "kube_pod_subnet"
	KubeServiceSubnetFactName        = "kube_service_subnet"
	KubeDnsDomainFactName            = "kube_dns_domain"
	KubernetesAuditFactName          = "kubernetes_audit"
	NodeportAddressFactName          = "nodeport_address"
	KubeServiceNodePortRangeFactName = "kube_service_node_port_range"

	DockerVersionFactName        = "docker_version"
	DockerMirrorRegistryFactName = "docker_mirror_registry"
	DockerRemoteApiFactName      = "docker_remote_api"
	ContainerdVersionFactName    = "containerd_version"
	ContainerRuntimeFactName     = "container_runtime"
	ContainerdStorageDirFactName = "containerd_storage_dir"
	DockerStorageDirFactName     = "docker_storage_dir"
	DockerSubnetFactName         = "docker_subnet"

	MasterScheduleTypeFactName = "master_schedule_type"
	LbModeFactName             = "lb_mode"
	LbKubeApiserverIpFactName  = "lb_kube_apiserver_ip"
	KubeApiserverPortFactName  = "lb_kube_apiserver_port"

	KubeProxyModeFactName           = "kube_proxy_mode"
	NetworkPluginFactName           = "network_plugin"
	CiliumVersionFactName           = "cilium_version"
	CiliumTunnelModeFactName        = "cilium_tunnel_mode"
	CiliumNativeRoutingCidrFactName = "cilium_native_routing_cidr"
	FlannelVersionFactName          = "flannel_version"
	FlannelBackendFactName          = "flannel_backend"
	CalicoVersionFactName           = "calico_version"
	CalicoIpv4poolIpIpFactName      = "calico_ipv4pool_ipip"
	NetworkInterfaceFactName        = "network_interface"
	NetworkCidrFactName             = "network_cidr"

	CgroupDriverFactName                 = "cgroup_driver"
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
	BinDirFactName                       = "bin_dir"
	BaseDirFactName                      = "base_dir"
	RegistryHostnameFactName             = "registry_hostname"
	RepoPortFactName                     = "repo_port"
	RegistryPortFactName                 = "registry_port"
	RegistryHostedPortFactName           = "registry_hosted_port"
	RegistryProtocolFactName             = "registry_protocol"
	CorednsImageFactName                 = "coredns_image"
	KubeadmTokenFactName                 = "kubeadm_token"

	HelmVersionFactName             = "helm_version"
	HelmV2VersionFactName           = "helm_v2_version"
	HelmV3VersionFactName           = "helm_v3_version"
	EtcdVersionFactName             = "etcd_version"
	EtcdDataDirFactName             = "etcd_data_dir"
	EtcdSnapshotCountFactName       = "etcd_snapshot_count"
	EtcdCompactionRetentionFactName = "etcd_compaction_retention"
	EtcdMaxRequestFactName          = "etcd_max_request_bytes"
	EtcdQuotaBackendFactName        = "etcd_quota_backend_bytes"
	CorednsVersionFactName          = "coredns_version"
	IngressControllerTypeFactName   = "ingress_controller_type"
	EnableNginxFactName             = "enable_nginx"
	NginxIngressVersionFactName     = "nginx_ingress_version"
	EnableTraefikFactName           = "enable_traefik"
	TraefikIngressVersionFactName   = "traefik_ingress_version"
	MetricsServerFactName           = "enable_metrics_server"
	MetricsServerVersionFactName    = "metrics_server_version"
	EnableDnsCacheFactName          = "enable_dns_cache"
	DnsCacheVersionFactName         = "dns_cache_version"
	NtpServerFactName               = "ntp_server"
	SupportGpuFactName              = "enable_gpu"
	EnableNpdFactName               = "enable_npd"
	EnableIstioFactName             = "enable_istio"

	ProvisionerNamespaceFactName = "provisioner_namespace"
	EnableNfsFactName            = "enable_nfs_provisioner"
	EnableGfsFactName            = "enable_gfs_provisioner"
	EnableCephBlockFactName      = "enable_external_ceph_block_provisioner"
	EnableCephFsFactName         = "enable_external_cephfs_provisioner"
	EnableCinderFactName         = "enable_cinder_provisioner"
	EnableVsphereFactName        = "enable_vsphere_provisioner"
	EnableOceanstorFactName      = "enable_oceanstor_provisioner"
	EnableRookFactName           = "enable_rook_provisioner"
)

var DefaultFacts = map[string]string{
	KubeVersionFactName:                  "v1.18.6",
	NodeNameRuleFactName:                 "hostname",
	ContainerRuntimeFactName:             "docker",
	CorednsImageFactName:                 "docker.io/kubeoperator/coredns:1.6.7",
	LbModeFactName:                       "internal",
	KubeApiserverPortFactName:            "8443",
	KubeDnsDomainFactName:                "cluster.local",
	KubePodSubnetFactName:                "10.244.0.0/18",
	KubeServiceSubnetFactName:            "10.244.64.0/18",
	DockerSubnetFactName:                 "172.17.0.1/16",
	KubeNetworkNodePrefixFactName:        "24",
	KubeMaxPodsFactName:                  "110",
	KubeProxyModeFactName:                "iptables",
	DnsCacheVersionFactName:              "1.17.0",
	NetworkPluginFactName:                "calico",
	CiliumVersionFactName:                "v1.9.5",
	CiliumTunnelModeFactName:             "vxlan",
	CiliumNativeRoutingCidrFactName:      "10.244.0.0/18",
	NodeportAddressFactName:              "",
	KubeServiceNodePortRangeFactName:     "30000-32767",
	KubeletRootDirFactName:               "/var/lib/kubelet",
	DockerMirrorRegistryFactName:         "enable",
	DockerRemoteApiFactName:              "disable",
	DockerStorageDirFactName:             "/var/lib/docker",
	ContainerdStorageDirFactName:         "/var/lib/containerd",
	BinDirFactName:                       "/usr/local/bin",
	BaseDirFactName:                      "/opt/kubeoperator",
	RegistryHostnameFactName:             "172.16.10.64",
	RepoPortFactName:                     "8081",
	RegistryPortFactName:                 "8082",
	RegistryHostedPortFactName:           "8083",
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
	FlannelBackendFactName:               "vxlan",
	HelmVersionFactName:                  "v3",
	EtcdVersionFactName:                  "v3.4.9",
	EtcdDataDirFactName:                  "/var/lib/etcd",
	EtcdSnapshotCountFactName:            "50000",
	EtcdCompactionRetentionFactName:      "1",
	EtcdMaxRequestFactName:               "10485760",
	EtcdQuotaBackendFactName:             "8589934592",
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
	NetworkCidrFactName:                  "",
	YumRepoFactName:                      "replace",
	NtpServerFactName:                    "ntp1.aliyun.com",
	MasterScheduleTypeFactName:           "enable",
	CgroupDriverFactName:                 "systemd",

	ComponentOptionFactName:       "component",
	IngressControllerTypeFactName: "nginx",
	EnableNginxFactName:           "enable",
	EnableTraefikFactName:         "disable",
	MetricsServerFactName:         "enable",
	SupportGpuFactName:            "disable",
	EnableDnsCacheFactName:        "enable",
	EnableNpdFactName:             "disable",
	EnableIstioFactName:           "disable",

	ProvisionerNamespaceFactName: "kube-system",
}
