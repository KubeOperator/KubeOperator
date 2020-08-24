package kubeconfig

import "github.com/KubeOperator/KubeOperator/pkg/util/ssh"

func ReadKubeConfigFile(client ssh.Interface) ([]byte, error) {
	return client.ReadFile("/root/.kube/config")
}
