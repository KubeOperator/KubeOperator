package encrypt

import (
	"fmt"
	"github.com/KubeOperator/KubeOperator/pkg/config"
	"testing"
)

func TestStringEncrypt(t *testing.T) {
	config.Init()
	p, err := StringEncrypt("kubeoperator@admin123")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(p)
}

func TestStringDecrypt(t *testing.T) {
	p, err := StringDecrypt("47zHCOqO84rdzGgxw5XPfgDEapoOMXbgJnryG32xp6Y=")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(p)
}
