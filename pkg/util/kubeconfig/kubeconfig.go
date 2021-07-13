package kubeconfig

import (
	"fmt"

	"github.com/KubeOperator/KubeOperator/pkg/util/ssh"
	"github.com/pkg/errors"
)

func ReadKubeConfigFile(client ssh.Interface, username string) ([]byte, error) {
	path := ""
	if username == "root" {
		path = "/root/.kube/config"
	} else {
		path = "/home/" + username + "/.kube/config"
	}
	result, err := client.ReadFile(path)
	if err != nil {
		return result, errors.Wrap(err, fmt.Sprintf("read file of %s failed: %v", path, err))
	}
	return result, err
}
