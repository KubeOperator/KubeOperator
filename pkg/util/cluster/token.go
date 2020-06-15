package cluster

import (
	"github.com/KubeOperator/KubeOperator/pkg/util/ssh"
	uuid "github.com/satori/go.uuid"
)

const (
	cmd = "kubectl -n kubernetes-system describe secret $(kubectl -n kubernetes-system get secret | grep ko-admin | awk '{print $1}') | grep token: | awk '{print $2}'"
)

func GetClusterToken(client ssh.Interface) (string, error) {
	buf, err := client.CombinedOutput(cmd)
	if err != nil {
		return "", err
	}
	return string(buf), nil
}

func GenerateKubeadmToken() string {
	return uuid.NewV4().String()
}
