package license

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os/exec"
	"path"
	"runtime"
	"strings"

	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/util/ssh"
)

type Response struct {
	Status  string      `json:"status"`
	License dto.License `json:"license"`
	Message string      `json:"message"`
}

func Parse(content string) (*Response, error) {
	if ssh.CheckIllegal(content) {
		return nil, errors.New("license contains invalid characters!")
	}
	fs, err := ioutil.ReadDir("/usr/local/bin")

	if err != nil {
		return nil, err
	}
	var validatorPath string
	for _, f := range fs {
		if strings.Contains(f.Name(), "validator") && strings.Contains(f.Name(), runtime.GOARCH) {
			validatorPath = path.Join("/usr/local/bin", f.Name())
		}
	}
	cmd := exec.Command(validatorPath, content)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	defer stdout.Close()
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	opBytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		return nil, err
	}
	var resp Response
	err = json.Unmarshal(opBytes, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
