package cluster

import (
	"github.com/KubeOperator/KubeOperator/pkg/util/ssh"
	uuid "github.com/satori/go.uuid"
	"strings"
)

const (
	cmd = "kubectl -n kube-system describe secret $(kubectl -n kube-system get secret | grep ko-admin | awk '{print $1}') | grep token: | awk '{print $2}'"
)

func GetClusterToken(client ssh.Interface) (string, error) {
	buf, err := client.CombinedOutput(cmd)
	if err != nil {
		return "", err
	}
	result := string(buf)
	result = strings.Replace(result, "\n", "", -1)
	return result, nil
}

func GenerateKubeadmToken() string {
	return uuid.NewV4().String()
}
