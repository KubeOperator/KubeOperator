package image

import (
	"fmt"
	"k8s.io/klog"
	"ko3-gin/pkg/cluster/adm/api"
	"ko3-gin/pkg/cluster/adm/constants"
	"ko3-gin/pkg/cluster/adm/util"
)

const extraHyperKubeNote = ` The "useHyperKubeImage" field will be removed from future kubeadm config versions and possibly ignored in future releases.`

// GetGenericImage generates and returns a platform agnostic image (backed by manifest list)
func GetGenericImage(prefix, image, tag string) string {
	return fmt.Sprintf("%s/%s:%s", prefix, image, tag)
}

// GetKubernetesImage generates and returns the image for the components managed in the Kubernetes main repository,
// including the control-plane components and kube-proxy. If specified, the HyperKube image will be used.
func GetKubernetesImage(image string, cfg *api.ClusterConfiguration) string {
	if image != constants.HyperKube {
		image = constants.HyperKube
	}
	repoPrefix := cfg.GetControlPlaneImageRepository()
	kubernetesImageTag := util.KubernetesVersionToImageTag(cfg.KubernetesVersion)
	return GetGenericImage(repoPrefix, image, kubernetesImageTag)
}

// GetDNSImage generates and returns the image for the DNS, that can be CoreDNS or kube-dns.
// Given that kube-dns uses 3 containers, an additional imageName parameter was added
func GetDNSImage(cfg *api.ClusterConfiguration, imageName string) string {
	// DNS uses default image repository by default
	dnsImageRepository := cfg.ImageRepository
	// unless an override is specified
	if cfg.DNS.ImageRepository != "" {
		dnsImageRepository = cfg.DNS.ImageRepository
	}
	// DNS uses an imageTag that corresponds to the DNS version.go matching the Kubernetes version.go
	dnsImageTag := constants.GetDNSVersion(cfg.DNS.Type)

	// unless an override is specified
	if cfg.DNS.ImageTag != "" {
		dnsImageTag = cfg.DNS.ImageTag
	}
	return GetGenericImage(dnsImageRepository, imageName, dnsImageTag)
}

// GetEtcdImage generates and returns the image for etcd
func GetEtcdImage(cfg *api.ClusterConfiguration) string {
	// Etcd uses default image repository by default
	etcdImageRepository := cfg.ImageRepository
	// unless an override is specified
	if cfg.Etcd.Local != nil && cfg.Etcd.Local.ImageRepository != "" {
		etcdImageRepository = cfg.Etcd.Local.ImageRepository
	}
	// Etcd uses an imageTag that corresponds to the etcd version.go matching the Kubernetes version.go
	etcdImageTag := constants.DefaultEtcdVersion
	etcdVersion, warning, err := constants.EtcdSupportedVersion(constants.SupportedEtcdVersion, cfg.KubernetesVersion)
	if err == nil {
		etcdImageTag = etcdVersion.String()
	}
	if warning != nil {
		klog.Warningln(warning)
	}
	// unless an override is specified
	if cfg.Etcd.Local != nil && cfg.Etcd.Local.ImageTag != "" {
		etcdImageTag = cfg.Etcd.Local.ImageTag
	}
	return GetGenericImage(etcdImageRepository, constants.Etcd, etcdImageTag)
}

// GetControlPlaneImages returns a list of container images kubeadm expects to use on a control plane node
func GetControlPlaneImages(cfg *api.ClusterConfiguration) []string {
	imgs := []string{}
	imgs = append(imgs, GetKubernetesImage(constants.KubeAPIServer, cfg))
	imgs = append(imgs, GetKubernetesImage(constants.KubeControllerManager, cfg))
	imgs = append(imgs, GetKubernetesImage(constants.KubeScheduler, cfg))
	imgs = append(imgs, GetKubernetesImage(constants.KubeProxy, cfg))

	// pause is not available on the ci image repository so use the default image repository.
	imgs = append(imgs, GetPauseImage(cfg))

	// if etcd is not external then add the image as it will be required
	if cfg.Etcd.Local != nil {
		imgs = append(imgs, GetEtcdImage(cfg))
	}

	// Append the appropriate DNS images
	if cfg.DNS.Type == api.CoreDNS {
		imgs = append(imgs, GetDNSImage(cfg, constants.CoreDNSImageName))
	} else {
		imgs = append(imgs, GetDNSImage(cfg, constants.KubeDNSKubeDNSImageName))
		imgs = append(imgs, GetDNSImage(cfg, constants.KubeDNSSidecarImageName))
		imgs = append(imgs, GetDNSImage(cfg, constants.KubeDNSDnsMasqNannyImageName))
	}

	return imgs
}

func GetPauseImage(cfg *api.ClusterConfiguration) string {
	return GetGenericImage(cfg.ImageRepository, "pause", constants.PauseVersion)
}
