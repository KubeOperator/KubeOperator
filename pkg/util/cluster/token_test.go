package cluster

import (
	"fmt"
	"github.com/KubeOperator/KubeOperator/pkg/util/ssh"
	"testing"
	"time"
)

func TestGetClusterToken(t *testing.T) {
	client, err := ssh.New(&ssh.Config{
		User:        "root",
		Host:        "172.16.10.184",
		Port:        22,
		Password:    "Calong@2015",
		PrivateKey:  nil,
		PassPhrase:  nil,
		DialTimeOut: 5 * time.Second,
		Retry:       3,
	})
	if err != nil {
		t.Error(err)
	}
	result, err := GetClusterToken(client)
	fmt.Println(result)
}
