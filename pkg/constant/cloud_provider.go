package constant

const (
	OpenStack                = "OpenStack"
	OpenStackImageName       = "kubeoperator_centos_7.6.1810"
	OpenStackImageDiskFormat = "qcow2"
	OpenStackImagePath       = "http://%s:8081/repository/oss-proxy/terraform/images/openstack/kubeoperator_centos_7.6.1810-1.qcow2"
	VSphere                  = "vSphere"
	VSphereImageName         = "kubeoperator_centos_7.6.1810"
	VSphereImageVMDkPath     = "http://%s:8081/repository/oss-proxy/terraform/images/vsphere/kubeoperator_centos_7.6.1810/kubeoperator_centos_7.6.1810-1.vmdk"
	VSphereImageOvfPath      = "http://%s:8081/repository/oss-proxy/terraform/images/vsphere/kubeoperator_centos_7.6.1810/kubeoperator_centos_7.6.1810.ovf"
	VSphereFolder            = "kubeoperator"
	ImageCredentialName      = "kubeoperator"
	OpenStackImageLocalPath  = "/opt/kubeoperator_centos_7.6.1810-1.qcow2"
	FusionCompute            = "FusionCompute"
	FusionComputeImageName   = "kubeoperator_centos_7.6.1810"
	FusionComputeOvfPath     = "http://%s:8081/repository/oss-proxy/terraform/images/fusioncompute/kubeoperator_centos_7.6.1810/kubeoperator_centos_7.6.1810.ovf"
	FusionComputeVhdPath     = "http://%s:8081/repository/oss-proxy/terraform/images/fusioncompute/kubeoperator_centos_7.6.1810/kubeoperator_centos_7.6.1810-vda.vhd"
	FusionComputeOvfName     = "kubeoperator_centos_7.6.1810.ovf"
	FusionComputeVhdName     = "kubeoperator_centos_7.6.1810-vda.vhd"
	FusionComputeOvfLocal    = "./kubeoperator_centos_7.6.1810.ovf"
	FusionComputeVhdLocal    = "./kubeoperator_centos_7.6.1810-vda.vhd"
)
