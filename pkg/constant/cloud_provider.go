package constant

const (
	OpenStack                = "OpenStack"
	OpenStackImageName       = "kubeoperator_centos_7.6.1810"
	OpenStackImageDiskFormat = "qcow2"
	OpenStackImagePath       = "http://%s:8081/repository/oss-proxy/terraform/images/openstack/kubeoperator_centos_7.6.1810-1.qcow2"
	VSphere                  = "vSphere"
	VSphereImageName         = "kubeoperator_centos_7.6.1810"
	VSphereImageVMDkPath     = "http://%s:8081/repository/oss-proxy/terraform/images/vsphere/kubeoperator_centos_7.6.1810/kubeoperator_centos_7.6.1810-1.vmdk.tar.gz"
	VSphereImageOvfPath      = "http://%s:8081/repository/oss-proxy/terraform/images/vsphere/kubeoperator_centos_7.6.1810/kubeoperator_centos_7.6.1810.ovf"
	VSphereFolder            = "kubeoperator"
	ImageDefaultPassword     = "KubeOperator@2019"
	ImageCredentialName      = "kubeoperator"
	ImageUserName            = "root"
	ImagePasswordType        = "password"
	VMDKGZLocalPath          = "/etc/ko/kubeoperator_centos_7.6.1810-1.vmdk.tar.gz"
	VMDKLocalPath            = "/etc/ko/kubeoperator_centos_7.6.1810-1.vmdk"
)
