package net

import (
	"fmt"
	"testing"
)

func TestPing(t *testing.T) {
	fmt.Println(Ping("172.16.10.184"))
}
