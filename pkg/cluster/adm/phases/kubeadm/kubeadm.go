package kubeadm

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"ko3-gin/pkg/cluster/adm/constants"
	"ko3-gin/pkg/cluster/adm/resource"
	"ko3-gin/pkg/cluster/util"
	"ko3-gin/pkg/util/ssh"
	"path"
)

const (
	kubeadmConfigFile  = "kubeadm/kubeadm-config.yaml"
	kubeadmKubeletConf = "/usr/lib/systemd/system/kubelet.service.d/10-kubeadm.conf"

	joinControlPlaneCmd = `kubeadm join {{.ControlPlaneEndpoint}} \
--node-name={{.NodeName}} --token={{.BootstrapToken}} \
--control-plane --certificate-key={{.CertificateKey}} \
--skip-phases=control-plane-join/mark-control-plane \
--discovery-token-unsafe-skip-ca-verification \
--ignore-preflight-errors=ImagePull \
--ignore-preflight-errors=Port-10250 \
--ignore-preflight-errors=FileContent--proc-sys-net-bridge-bridge-nf-call-iptables \
--ignore-preflight-errors=DirAvailable--etc-kubernetes-manifests
`
	joinNodeCmd = `kubeadm join {{.ControlPlaneEndpoint}} \
--node-name={{.NodeName}} \
--token={{.BootstrapToken}} \
--discovery-token-unsafe-skip-ca-verification \
--ignore-preflight-errors=ImagePull \
--ignore-preflight-errors=Port-10250 \
--ignore-preflight-errors=FileContent--proc-sys-net-bridge-bridge-nf-call-iptables
`
)

func Install(s ssh.Interface) error {
	dstFile, err := resource.Kubeadm.CopyToNodeWithDefault(s)
	if err != nil {
		return err
	}

	cmd := "tar xvaf %s -C %s "
	_, stderr, exit, err := s.Exec("tar", "xvaf", dstFile, "-C", constants.BinDir)
	if err != nil || exit != 0 {
		return fmt.Errorf("exec %q failed:exit %d:stderr %s:error %s", cmd, exit, stderr, err)
	}

	data, err := util.ParseFile(path.Join(constants.ConfDir, "kubeadm/10-kubeadm.conf"), nil)
	if err != nil {
		return err
	}
	err = s.WriteFile(bytes.NewReader(data), kubeadmKubeletConf)
	if err != nil {
		return errors.Wrapf(err, "write %s error", kubeadmKubeletConf)
	}

	return nil
}

type InitOption struct {
	KubeadmConfigFileName string
	NodeName              string
	BootstrapToken        string
	CertificateKey        string

	ETCDImageTag         string
	CoreDNSImageTag      string
	KubernetesVersion    string
	ControlPlaneEndpoint string

	DNSDomain             string
	ServiceSubnet         string
	NodeCIDRMaskSize      int32
	ClusterCIDR           string
	ServiceClusterIPRange string
	CertSANs              []string

	APIServerExtraArgs         map[string]string
	ControllerManagerExtraArgs map[string]string
	SchedulerExtraArgs         map[string]string

	ImageRepository string
	ClusterName     string

	KubeProxyMode string
}

func Init(s ssh.Interface, option *InitOption, extraCmd string) error {
	configData, err := util.ParseFile(path.Join(constants.ConfDir, kubeadmConfigFile), option)
	if err != nil {
		return errors.Wrap(err, "parse kubeadm config file error")
	}
	err = s.WriteFile(bytes.NewReader(configData), option.KubeadmConfigFileName)
	if err != nil {
		return err
	}

	cmd := fmt.Sprintf("kubeadm init phase %s --config=%s",
		extraCmd, option.KubeadmConfigFileName)
	_, stderr, exit, err := s.Exec(cmd)
	if err != nil || exit != 0 {
		return fmt.Errorf("exec %q failed:exit %d:stderr %s:error %s", cmd, exit, stderr, err)
	}

	return nil
}

type JoinControlPlaneOption struct {
	NodeName             string
	BootstrapToken       string
	CertificateKey       string
	ControlPlaneEndpoint string
}

func JoinControlPlane(s ssh.Interface, option *JoinControlPlaneOption) error {
	cmd, err := util.ParseString(joinControlPlaneCmd, option)
	if err != nil {
		return errors.Wrap(err, "parse joinControlePlaneCmd error")
	}
	_, stderr, exit, err := s.Exec(string(cmd))
	if err != nil || exit != 0 {
		return fmt.Errorf("exec %q failed:exit %d:stderr %s:error %s", cmd, exit, stderr, err)
	}
	return nil
}

type JoinNodeOption struct {
	NodeName             string
	BootstrapToken       string
	ControlPlaneEndpoint string
}

func JoinNode(s ssh.Interface, option *JoinNodeOption) error {
	cmd, err := util.ParseString(joinNodeCmd, option)
	if err != nil {
		return errors.Wrap(err, "parse joinNodeCmd error")
	}
	_, stderr, exit, err := s.Exec(string(cmd))
	if err != nil || exit != 0 {
		_, _, _, _ = s.Exec("kubeadm reset -f")
		return fmt.Errorf("exec %q failed:exit %d:stderr %s:error %s", cmd, exit, stderr, err)
	}
	return nil
}
