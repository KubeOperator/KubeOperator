package encrypt

import (
	"fmt"
	"github.com/KubeOperator/KubeOperator/pkg/config"
	"testing"
)

func TestStringEncrypt(t *testing.T) {
	config.Init()
	p, err := StringEncrypt("")
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
