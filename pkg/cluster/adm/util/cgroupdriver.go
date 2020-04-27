package util

import (
	"ko3-gin/pkg/util/ssh"
	"strings"

	"github.com/pkg/errors"
)

const (
	CgroupDriverSystemd  = "systemd"
	CgroupDriverCgroupfs = "cgroupfs"
)

func GetCgroupDriverDocker(ssh *ssh.SSH) (string, error) {
	driver, err := callDockerInfo(ssh)
	if err != nil {
		return "", err
	}
	return strings.TrimSuffix(driver, "\n"), nil
}

func callDockerInfo(ssh *ssh.SSH) (string, error) {
	out, err := ssh.CombinedOutput("docker", "info", "-f", "{{.CgroupDriver}}")
	if err != nil {
		return "", errors.Wrap(err, "cannot execute 'docker info -f {{.CgroupDriver}}'")
	}
	return string(out), nil
}
