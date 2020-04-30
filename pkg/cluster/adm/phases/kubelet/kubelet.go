package kubelet

import (
	"bytes"
	"fmt"
	"ko3-gin/pkg/cluster/adm/constants"
	"ko3-gin/pkg/cluster/adm/resource"
	"ko3-gin/pkg/cluster/util"
	"ko3-gin/pkg/util/ssh"
	"path"
	"strings"
)

type Option struct {
	Version   string
	ExtraArgs map[string]string
}

func Install(s ssh.Interface, option *Option) error {
	dstFile, err := resource.KubernetesNode.CopyToNode(s, option.Version)
	if err != nil {
		return err
	}

	cmd := "tar xvaf %s -C %s --strip-components=3"
	_, stderr, exit, err := s.Exec("tar", "xvaf", dstFile, "-C", constants.BinDir, "--strip-components=3")
	if err != nil {
		return fmt.Errorf("exec %q failed:exit %d:stderr %s:error %s", cmd, exit, stderr, err)
	}

	var args []string
	for k, v := range option.ExtraArgs {
		args = append(args, fmt.Sprintf(`--%s="%s"`, k, v))
	}
	err = s.WriteFile(strings.NewReader(fmt.Sprintf("KUBELET_EXTRA_ARGS=%s", strings.Join(args, " "))), "/etc/sysconfig/kubelet")
	if err != nil {
		return err
	}

	serviceData, err := util.ParseFile(path.Join(constants.ConfDir, "kubelet/kubelet.service"), nil)
	if err != nil {
		return err
	}

	ss := &util.SystemCtl{Name: "kubelet", SSH: s}
	err = ss.Deploy(bytes.NewReader(serviceData))
	if err != nil {
		return err
	}

	err = ss.Start()
	if err != nil {
		return err
	}

	return nil
}
