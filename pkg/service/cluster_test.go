package service

import (
	"fmt"
	"github.com/KubeOperator/KubeOperator/pkg/config"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/service/dto"
	"github.com/spf13/viper"
	"log"
	"testing"
)

func Init() {
	config.Init()
	phase := db.InitDBPhase{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetInt("db.port"),
		Name:     viper.GetString("db.name"),
		User:     viper.GetString("db.user"),
		Password: viper.GetString("db.password"),
	}
	err := phase.Init()
	if err != nil {
		log.Fatalf("can not init db,%s", err)
	}
}

func TestClusterService_Create(t *testing.T) {
	Init()
	service := NewClusterService()
	if err := service.Create(dto.ClusterCreate{
		Name:                 "sdf",
		Version:              "1.18.0",
		NetworkType:          "flannel",
		RuntimeType:          "docker",
		DockerStorageDIr:     "/var/lib/docker",
		ContainerdStorageDIr: "",
		AppDomain:            "apps.sdf.com",
		Nodes: []dto.NodeCreate{
			{
				Role:     "master",
				HostName: "node-1",
			},
			{
				Role:     "worker",
				HostName: "node-2",
			},
			{
				Role:     "worker",
				HostName: "node-3",
			},
		},
	}); err != nil {
		t.Error(err)
	}
}

func TestClusterService_GetStatus(t *testing.T) {
	Init()
	service := NewClusterService()
	status, err := service.GetStatus("sdf")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(status)
}

func TestClusterService_Delete(t *testing.T) {
	Init()
	service := NewClusterService()
	err := service.Delete("sdf")
	if err != nil {
		t.Error(err)
	}
}
