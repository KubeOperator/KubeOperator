package cluster

import (
	"fmt"
	"testing"
	"time"

	"github.com/KubeOperator/KubeOperator/pkg/logger"
	"github.com/KubeOperator/KubeOperator/pkg/util/ssh"
)

func TestGetClusterToken(t *testing.T) {
	client, err := ssh.New(&ssh.Config{
		User:        "root",
		Host:        "172.16.10.210",
		Port:        22,
		Password:    "KubeOperator@2019",
		PrivateKey:  nil,
		PassPhrase:  nil,
		DialTimeOut: 5 * time.Second,
		Retry:       3,
	})
	if err != nil {
		t.Error(err)
	}
	by, err := client.CombinedOutput("kubectl", "get sa -A | grep default")
	if err != nil {
		logger.Log.Fatal(err)
	}
	fmt.Println(by)

	//result, err := GetClusterToken(client)
	//fmt.Println(result, err.Error())
}
