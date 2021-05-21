package ssh

import (
	"fmt"
	"io/ioutil"
	"testing"
	"time"

	"github.com/KubeOperator/KubeOperator/pkg/logger"
)

func TestSSHClient(t *testing.T) {
	key, err := ioutil.ReadFile("/Users/shenchenyang/Desktop/aa.key")
	if err != nil {
		logger.Log.Fatal(err)
	}

	client, err := New(&Config{
		User:        "root",
		Host:        "172.16.10.210",
		Port:        22,
		PrivateKey:  key,
		PassPhrase:  nil,
		DialTimeOut: 5 * time.Second,
		Retry:       3,
	})
	if err != nil {
		logger.Log.Fatal(err)
	}
	if err := client.Ping(); err != nil {
		logger.Log.Fatal(err)
	}
	bs, err := client.CombinedOutput("ps", "aux")
	if err != nil {
		logger.Log.Fatal(err)
	}
	fmt.Println(string(bs))
}
