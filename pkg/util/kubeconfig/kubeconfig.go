package kubeconfig

import (
	"fmt"

	"github.com/KubeOperator/KubeOperator/pkg/util/ssh"
	"github.com/pkg/errors"
)

func ReadKubeConfigFile(client ssh.Interface) ([]byte, error) {
	result, err := client.ReadFile("/root/.kube/config")
	if err != nil {
		return result, errors.Wrap(err, fmt.Sprintf("read file of /root/.kube/config failed: %v", err))
	}
	return result, err
}
