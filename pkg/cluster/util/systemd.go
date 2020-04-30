package util

import (
	"fmt"
	"io"
	"ko3-gin/pkg/cluster/adm/constants"
	"ko3-gin/pkg/util/ssh"
	"path"
)

type SystemCtl struct {
	Name string
	SSH  ssh.Interface
}

func (s *SystemCtl) Deploy(data io.Reader) error {
	unitFilePath := path.Join(constants.DefaultSystemdUnitFilePath, fmt.Sprintf("%s.service", s.Name))
	if err := s.SSH.WriteFile(data, unitFilePath); err != nil {
		return err
	}

	if _, stderr, exit, err := s.SSH.Exec("systemctl", "-f", "enable", unitFilePath); err != nil || exit != 0 {
		return fmt.Errorf("exec %q failed:exit %d:stderr %s:error %s", "", exit, stderr, err)
	}

	if _, stderr, exit, err := s.SSH.Exec("systemctl", "daemon-reload"); err != nil || exit != 0 {
		return fmt.Errorf("exec %q failed:exit %d:stderr %s:error %s", "", exit, stderr, err)
	}

	return nil
}

func (s *SystemCtl) Start() error {
	unitName := fmt.Sprintf("%s.service", s.Name)

	if _, stderr, exit, err := s.SSH.Exec("systemctl", "restart", unitName); err != nil || exit != 0 {
		jStdout, _, jExit, jErr := s.SSH.Exec("journalctl", "--unit", unitName, "-n10", "--no-pager")
		if jErr != nil || jExit != 0 {
			return fmt.Errorf("exec %q:error %s", "", err)
		}
		fmt.Printf("log:\n%s", jStdout)

		return fmt.Errorf("Exec %s failed:exit %d:stderr %s:error %s:log:\n%s", "", exit, stderr, err, jStdout)
	}

	return nil
}
