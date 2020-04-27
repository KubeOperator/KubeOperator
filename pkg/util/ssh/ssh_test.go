package ssh

import (
	"fmt"
	"testing"
	"time"
)

func TestSSHClient(t *testing.T) {
	client, err := New(&Config{
		User:        "root",
		Host:        "119.28.214.236",
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
	if err := client.Ping(); err != nil {
		t.Error(err)
	}
	bs, err := client.CombinedOutput("ps", "aux")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(bs))
}
