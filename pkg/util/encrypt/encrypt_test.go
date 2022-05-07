package encrypt

import (
	"fmt"
	"testing"

	"github.com/KubeOperator/KubeOperator/pkg/config"
)

func TestStringEncrypt(t *testing.T) {
	config.Init()
	p, err := StringEncrypt("kubepi")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(p)
}

func TestStringDecrypt(t *testing.T) {
	p, err := StringDecrypt("")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(p)
}
