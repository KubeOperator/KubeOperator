package cluster

import (
	"strings"
	"time"

	"github.com/KubeOperator/KubeOperator/pkg/logger"
	"github.com/KubeOperator/KubeOperator/pkg/util/ssh"
	uuid "github.com/satori/go.uuid"
	"k8s.io/apimachinery/pkg/util/wait"
)

const (
	cmd = "/usr/local/bin/kubectl get sa -A | grep ko-admin &> /dev/null && /usr/local/bin/kubectl -n kube-system describe secret $(kubectl -n kube-system get secret | grep ko-admin | awk '{print $1}') | grep token: | awk '{print $2}'"
)

var log = logger.Default

func GetClusterToken(client ssh.Interface) (string, error) {
	result := ""
	if err := wait.Poll(5*time.Second, 1*time.Minute, func() (done bool, err error) {
		buf, err := client.CombinedOutput(cmd)
		if err != nil || len(buf) == 0 {
			log.Error("can not get kubernetes token ,retry after 5 second")
			return false, nil
		}
		result = string(buf)
		result = strings.Replace(result, "\n", "", -1)
		return true, nil
	}); err != nil {
		return "", err
	}
	return result, nil
}

func GenerateKubeadmToken() string {
	return uuid.NewV4().String()
}
