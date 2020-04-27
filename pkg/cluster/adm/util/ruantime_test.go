package util

import (
	"fmt"
	"ko3-gin/pkg/util/ssh"
	"log"
	"testing"
	"time"
)

func GetRuntime() ContainerRuntime {
	client, _ := ssh.New(&ssh.Config{
		User:        "root",
		Host:        "119.28.214.236",
		Port:        22,
		Password:    "Calong@2015",
		DialTimeOut: 5 * time.Second,
		Retry:       3,
	})
	runtime, err := NewContainerRuntime(client)
	if err != nil {
		log.Fatal(err)
	}
	return runtime
}

func TestDockerRuntime_IsDocker(t *testing.T) {
	runtime := GetRuntime()
	fmt.Println(runtime.IsDocker())
}

func TestDockerRuntime_IsRunning(t *testing.T) {
	runtime := GetRuntime()
	fmt.Println(runtime.IsRunning())
}

func TestDockerRuntime_ImageExists(t *testing.T) {
	runtime := GetRuntime()
	b, _ := runtime.ImageExists("centos:7")
	fmt.Println(b)
}
