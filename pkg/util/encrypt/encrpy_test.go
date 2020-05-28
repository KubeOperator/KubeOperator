package encrypt

import (
	"fmt"
	"testing"
)

func TestEncrypt(t *testing.T) {
	password, err := StringEncrypt("Calong@2015")
	if err == nil {
		fmt.Println(password)
		password2, _ := StringDecrypt(password)
		fmt.Println(password2)
	}
}
