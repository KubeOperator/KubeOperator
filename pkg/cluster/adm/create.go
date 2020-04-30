package adm

import (
	"bytes"
	"fmt"
	"github.com/google/martian/log"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/util/wait"
	"ko3-gin/pkg/cluster/adm/constants"
	"ko3-gin/pkg/cluster/adm/images"
	"ko3-gin/pkg/cluster/adm/phases/docker"
	"ko3-gin/pkg/cluster/adm/phases/kubeadm"
	"ko3-gin/pkg/cluster/adm/phases/kubeconfig"
	"ko3-gin/pkg/cluster/adm/phases/kubelet"
	"ko3-gin/pkg/cluster/util"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
)

const (
	sysctlFile       = "/etc/sysctl.conf"
	sysctlCustomFile = "/etc/sysctl.d/ko.conf"
	moduleFile       = "/etc/modules-load.d/ko.conf"
)

func (ca *ClusterAdm) EnsureCopyFiles(c *Cluster) error {
	return nil
}

func (ca *ClusterAdm) EnsureKernelModule(c *Cluster) error {
	modules := []string{"iptable_nat"}
	var data bytes.Buffer
	for _, machine := range c.Spec.Machines {
		s := c.SSH[machine.IP]

		for _, m := range modules {
			_, err := s.CombinedOutput(fmt.Sprintf("modprobe %s", m))
			if err != nil {
				return errors.Wrap(err, machine.IP)
			}
			data.WriteString(m + "\n")
		}
		err := s.WriteFile(strings.NewReader(data.String()), moduleFile)
		if err != nil {
			return errors.Wrap(err, machine.IP)
		}
	}

	return nil
}

func (ca *ClusterAdm) EnsureSysctl(c *Cluster) error {
	for _, machine := range c.Spec.Machines {
		s := c.SSH[machine.IP]

		_, err := s.CombinedOutput(util.SetFileContent(sysctlFile, "^net.ipv4.ip_forward.*", "net.ipv4.ip_forward = 1"))
		if err != nil {
			return errors.Wrap(err, machine.IP)
		}

		_, err = s.CombinedOutput(util.SetFileContent(sysctlFile, "^net.bridge.bridge-nf-call-iptables.*", "net.bridge.bridge-nf-call-iptables = 1"))
		if err != nil {
			return errors.Wrap(err, machine.IP)
		}

		f, err := os.Open(path.Join(constants.ConfDir, "sysctl.conf"))
		if err == nil {
			err = s.WriteFile(f, sysctlCustomFile)
			if err != nil {
				return err
			}
		}

		_, err = s.CombinedOutput("sysctl --system")
		if err != nil {
			return errors.Wrap(err, machine.IP)
		}
	}

	return nil
}

func (ca *ClusterAdm) EnsureDisableSwap(c *Cluster) error {
	for _, machine := range c.Spec.Machines {
		s := c.SSH[machine.IP]

		_, err := s.CombinedOutput("swapoff -a && sed -i 's/^[^#]*swap/#&/' /etc/fstab")
		if err != nil {
			return errors.Wrap(err, machine.IP)
		}
	}

	return nil
}

func (ca *ClusterAdm) EnsureClusterComplete(cluster *Cluster) error {
	return nil
}

func (ca *ClusterAdm) EnsureDocker(c *Cluster) error {
	insecureRegistries := fmt.Sprintf(`"%s"`, c.Registry.Domain)
	option := &docker.Option{
		InsecureRegistries: insecureRegistries,
		RegistryDomain:     c.Registry.Domain,
		ExtraArgs:          c.Spec.DockerExtraArgs,
	}
	for _, machine := range c.Spec.Machines {
		err := docker.Install(c.SSH[machine.IP], option)
		if err != nil {
			return errors.Wrap(err, machine.IP)
		}
	}
	return nil
}

func (ca *ClusterAdm) EnsureKubeadm(c *Cluster) error {
	for _, machine := range c.Spec.Machines {
		err := kubeadm.Install(c.SSH[machine.IP])
		if err != nil {
			return errors.Wrap(err, machine.IP)
		}
	}

	return nil
}

func (ca *ClusterAdm) EnsureKubelet(c *Cluster) error {
	option := &kubelet.Option{
		Version:   c.Spec.Version,
		ExtraArgs: c.Spec.KubeletExtraArgs,
	}
	for _, machine := range c.Spec.Machines {
		err := kubelet.Install(c.SSH[machine.IP], option)
		if err != nil {
			return errors.Wrap(err, machine.IP)
		}
	}

	return nil
}

func (ca *ClusterAdm) EnsurePrepareForControlplane(c *Cluster) error {
	for _, machine := range c.Spec.Machines {
		tokenData := fmt.Sprintf(tokenFileTemplate, *c.Credential.Token)
		err := c.SSH[machine.IP].WriteFile(strings.NewReader(tokenData), constants.TokenFile)
		if err != nil {
			return errors.Wrap(err, machine.IP)
		}
	}
	return nil
}

func (ca *ClusterAdm) EnsureKubeadmInitKubeletStartPhase(c *Cluster) error {
	return kubeadm.Init(c.SSH[c.Spec.Machines[0].IP], getKubeadmInitOption(c),
		fmt.Sprintf("kubelet-start --node-name=%s", c.Spec.Machines[0].IP))
}

func (ca *ClusterAdm) EnsureKubeadmInitCertsPhase(c *Cluster) error {
	return kubeadm.Init(c.SSH[c.Spec.Machines[0].IP], getKubeadmInitOption(c), "certs all")
}

func (ca *ClusterAdm) EnsureKubeadmInitKubeConfigPhase(c *Cluster) error {
	return kubeadm.Init(c.SSH[c.Spec.Machines[0].IP], getKubeadmInitOption(c), "kubeconfig all")
}

func (ca *ClusterAdm) EnsureKubeadmInitControlPlanePhase(c *Cluster) error {
	return kubeadm.Init(c.SSH[c.Spec.Machines[0].IP], getKubeadmInitOption(c), "control-plane all")
}

func (ca *ClusterAdm) EnsureKubeadmInitEtcdPhase(c *Cluster) error {
	return kubeadm.Init(c.SSH[c.Spec.Machines[0].IP], getKubeadmInitOption(c), "etcd local")
}

func (ca *ClusterAdm) EnsureKubeadmInitUploadConfigPhase(c *Cluster) error {
	return kubeadm.Init(c.SSH[c.Spec.Machines[0].IP], getKubeadmInitOption(c), "upload-config all ")
}

func (ca *ClusterAdm) EnsureKubeadmInitUploadCertsPhase(c *Cluster) error {
	return kubeadm.Init(c.SSH[c.Spec.Machines[0].IP], getKubeadmInitOption(c), "upload-certs --upload-certs")
}

func (ca *ClusterAdm) EnsureKubeadmInitBootstrapTokenPhase(c *Cluster) error {
	return kubeadm.Init(c.SSH[c.Spec.Machines[0].IP], getKubeadmInitOption(c), "bootstrap-token")
}

func (ca *ClusterAdm) EnsureKubeadmInitAddonPhase(c *Cluster) error {
	return kubeadm.Init(c.SSH[c.Spec.Machines[0].IP], getKubeadmInitOption(c), "addon all")
}

func (ca *ClusterAdm) EnsureStoreCredential(c *Cluster) error {
	data, err := c.SSH[c.Spec.Machines[0].IP].ReadFile(constants.CACertName)
	if err != nil {
		return errors.Wrapf(err, "read %s error", constants.CACertName)
	}
	c.Credential.CACert = data

	data, err = c.SSH[c.Spec.Machines[0].IP].ReadFile(constants.CAKeyName)
	if err != nil {
		return errors.Wrapf(err, "read %s error", constants.CAKeyName)
	}
	c.Credential.CAKey = data

	data, err = c.SSH[c.Spec.Machines[0].IP].ReadFile(constants.EtcdCACertName)
	if err != nil {
		return errors.Wrapf(err, "read %s error", constants.EtcdCACertName)
	}
	c.Credential.ETCDCACert = data

	data, err = c.SSH[c.Spec.Machines[0].IP].ReadFile(constants.EtcdCAKeyName)
	if err != nil {
		return errors.Wrapf(err, "read %s error", constants.EtcdCAKeyName)
	}
	c.Credential.ETCDCAKey = data

	data, err = c.SSH[c.Spec.Machines[0].IP].ReadFile(constants.APIServerEtcdClientCertName)
	if err != nil {
		return errors.Wrapf(err, "read %s error", constants.APIServerEtcdClientCertName)
	}
	c.Credential.ETCDAPIClientCert = data

	data, err = c.SSH[c.Spec.Machines[0].IP].ReadFile(constants.APIServerEtcdClientKeyName)
	if err != nil {
		return errors.Wrapf(err, "read %s error", constants.APIServerEtcdClientKeyName)
	}
	c.Credential.ETCDAPIClientKey = data

	return nil
}

func (ca *ClusterAdm) EnsureKubeconfig(c *Cluster) error {
	for _, machine := range c.Spec.Machines {
		option := &kubeconfig.Option{
			MasterEndpoint: "https://127.0.0.1:6443",
			ClusterName:    c.Name,
			CACert:         c.Credential.CACert,
			Token:          *c.Credential.Token,
		}
		err := kubeconfig.Install(c.SSH[machine.IP], option)
		if err != nil {
			return errors.Wrap(err, machine.IP)
		}
	}

	return nil
}

func (ca *ClusterAdm) EnsureKubeadmInitWaitControlPlanePhase(c *Cluster) error {
	start := time.Now()

	return wait.PollImmediate(5*time.Second, 5*time.Minute, func() (bool, error) {
		healthStatus := 0
		clientset, err := c.Clientset()
		if err != nil {
			return false, nil
		}
		clientset.Discovery().RESTClient().Get().AbsPath("/healthz").Do().StatusCode(&healthStatus)
		if healthStatus != http.StatusOK {
			return false, nil
		}

		log.Infof("All control plane components are healthy after %f seconds\n", time.Since(start).Seconds())
		return true, nil
	})
}

func (ca *ClusterAdm) EnsureJoinControlePlane(c *Cluster) error {
	option := &kubeadm.JoinControlPlaneOption{
		BootstrapToken:       *c.Credential.BootstrapToken,
		CertificateKey:       *c.Credential.CertificateKey,
		ControlPlaneEndpoint: fmt.Sprintf("%s:6443", c.Spec.Machines[0].IP),
	}
	for _, machine := range c.Spec.Machines[1:] {
		option.NodeName = machine.IP
		err := kubeadm.JoinControlPlane(c.SSH[machine.IP], option)
		if err != nil {
			return errors.Wrap(err, machine.IP)
		}
	}
	return nil
}

func getKubeadmInitOption(c *Cluster) *kubeadm.InitOption {
	controlPlaneEndpoint := fmt.Sprintf("%s:6443", c.Spec.Machines[0].IP)
	kubeProxyMode := "iptables"
	return &kubeadm.InitOption{
		KubeadmConfigFileName: constants.KubectlConfigFile,
		NodeName:              c.Spec.Machines[0].IP,
		BootstrapToken:        *c.Credential.BootstrapToken,
		CertificateKey:        *c.Credential.CertificateKey,

		ETCDImageTag:         images.Get().ETCD.Tag,
		CoreDNSImageTag:      images.Get().CoreDNS.Tag,
		KubernetesVersion:    c.Spec.Version,
		ControlPlaneEndpoint: controlPlaneEndpoint,

		DNSDomain:             c.Spec.DNSDomain,
		ServiceSubnet:         c.Spec.ServiceCIDR,
		ClusterCIDR:           c.Spec.ClusterCIDR,
		ServiceClusterIPRange: c.Spec.ServiceCIDR,
		CertSANs:              []string{},

		APIServerExtraArgs:         c.Spec.APIServerExtraArgs,
		ControllerManagerExtraArgs: c.Spec.ControllerManagerExtraArgs,
		SchedulerExtraArgs:         c.Spec.SchedulerExtraArgs,

		ImageRepository: c.Registry.Prefix,
		ClusterName:     c.Name,

		KubeProxyMode: kubeProxyMode,
	}
}
