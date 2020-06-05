package cluster

import "github.com/KubeOperator/KubeOperator/pkg/util/ssh"

const (
	cmd = "kubectl -n kube-system describe secret $(kubectl -n kube-system get secret | grep ko-admin | awk '{print $1}') | grep token: | awk '{print $2}'"
)

func GetClusterToken(client ssh.Interface) (string, error) {
	buf, err := client.CombinedOutput(cmd)
	if err != nil {
		return "", err
	}
	return string(buf), nil
}
