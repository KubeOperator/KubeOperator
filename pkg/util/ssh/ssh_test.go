package ssh

import (
	"fmt"
	"io/ioutil"
	"testing"
	"time"
)

func TestSSHClient(t *testing.T) {
	key, err := ioutil.ReadFile("/Users/shenchenyang/Desktop/aa.key")
	if err != nil {
		log.Fatal(err)
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
		log.Fatal(err)
	}
	if err := client.Ping(); err != nil {
		log.Fatal(err)
	}
	bs, err := client.CombinedOutput("ps", "aux")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(bs))
}
