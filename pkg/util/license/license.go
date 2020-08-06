package license

import (
	"encoding/json"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"io/ioutil"
	"os/exec"
)

type Response struct {
	Status  string      `json:"status"`
	License dto.License `json:"license"`
	Message string      `json:"message"`
}

func Parse(content string) (*Response, error) {
	cmd := exec.Command("/Users/shenchenyang/go/bin/validator_darwin_amd64", content)
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
