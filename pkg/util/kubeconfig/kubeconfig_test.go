package kubeconfig

import (
	"fmt"
	"github.com/KubeOperator/KubeOperator/pkg/util/ssh"
	"log"
	"testing"
)

func TestReadKubeConfigFile(t *testing.T) {
	s, err := ssh.New(&ssh.Config{
		User:        "root",
		Host:        "172.16.10.111",
		Port:        22,
		Password:    "KubeOperator@2019",
		Retry:       5,
	})
	if err != nil {
		log.Fatal(err.Error())
	}
	bs, err := ReadKubeConfigFile(s)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println(string(bs))
}
