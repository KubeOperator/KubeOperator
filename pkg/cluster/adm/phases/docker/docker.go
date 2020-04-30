package docker

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"ko3-gin/pkg/cluster/adm/constants"
	"ko3-gin/pkg/cluster/util"
	"ko3-gin/pkg/util/ssh"
	"path"
	"strings"
)

type Option struct {
	InsecureRegistries string
	RegistryDomain     string
	Options            string
	ExtraArgs          map[string]string
}

const (
	dockerDaemonFile = "/etc/docker/daemon.json"
)

func Install(s ssh.Interface, option *Option) error {
	// 下载文件
	//dstFile, err := res.Docker.CopyToNodeWithDefault(s)
	//if err != nil {
	//	return err
	//}

	cmd := fmt.Sprintf("tar xvaf %s -C %s --strip-components=1", "", constants.BinDir)
	_, stderr, exit, err := s.Exec(cmd)
	if err != nil {
		return fmt.Errorf("exec %q failed:exit %d:stderr %s:error %s", cmd, exit, stderr, err)
	}

	var args []string
	for k, v := range option.ExtraArgs {
		args = append(args, fmt.Sprintf(`--%s="%s"`, k, v))
	}
	err = s.WriteFile(strings.NewReader(fmt.Sprintf("DOCKER_EXTRA_ARGS=%s", strings.Join(args, " "))), "/etc/sysconfig/docker")
	if err != nil {
		return err
	}

	data, err := util.ParseFile(path.Join(constants.ConfDir, "docker/daemon.json"), option)
	if err != nil {
		return err
	}
	err = s.WriteFile(bytes.NewReader(data), dockerDaemonFile)
	if err != nil {
		return errors.Wrapf(err, "write %s error", dockerDaemonFile)
	}

	data, err = util.ParseFile(path.Join(constants.ConfDir, "docker/docker.service"), option)
	if err != nil {
		return err
	}
	ss := &util.SystemCtl{Name: "docker", SSH: s}
	err = ss.Deploy(bytes.NewReader(data))
	if err != nil {
		return err
	}

	err = ss.Start()
	if err != nil {
		return err
	}
	return nil
}
