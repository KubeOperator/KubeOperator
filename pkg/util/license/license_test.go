package license

import (
	"fmt"
	"io/ioutil"
	"log"
	"testing"
)

func TestParse(t *testing.T) {
	fs, err := ioutil.ReadDir("/usr/local/bin")
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range fs {
		fmt.Println(f.Name())
	}
}
