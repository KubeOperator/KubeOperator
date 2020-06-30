package service

import (
	"fmt"
	"github.com/KubeOperator/KubeOperator/pkg/config"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
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

func TestClusterStorageProvisionerService_CreateStorageProvisioner(t *testing.T) {
	Init()
	service := NewClusterStorageProvisionerService()
	p, err := service.CreateStorageProvisioner("test", dto.ClusterStorageProvisionerCreation{
		Name: "nfs-1",
		Type: "nfs",
		Vars: map[string]interface{}{
			"server": "172.16.10.69",
			"path":   "/expose",
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(p.Name)

}
