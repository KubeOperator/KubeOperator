package license

import (
	"encoding/json"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"io/ioutil"
	"os/exec"
	"path"
	"runtime"
	"strings"
)

func Parse(content string) (*dto.License, error) {
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
	var resp dto.License
	err = json.Unmarshal(opBytes, &resp)
	if err != nil {
		return nil, err
	}
	if err := cmd.Wait(); err != nil {
		return nil, err
	}
	return &resp, nil
}
