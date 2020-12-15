package net

import (
	"fmt"
	"os"
	"testing"
)

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

func TestPing(t *testing.T) {

}

func TestTcpPing(t *testing.T) {
	if err:=TcpPing("172.16.10.222:8443",true);err!=nil{
		t.Fatal(err)
	}

}
