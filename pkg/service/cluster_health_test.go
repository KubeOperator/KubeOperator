package service

import (
	"fmt"
	"github.com/KubeOperator/KubeOperator/pkg/config"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/spf13/viper"
	"testing"
)

func TestClusterHealthService_HealthCheck(t *testing.T) {
	config.Init()
	dbi := db.InitDBPhase{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetInt("db.port"),
		Name:     viper.GetString("db.name"),
		User:     viper.GetString("db.user"),
		Password: viper.GetString("db.password"),
	}
	err := dbi.Init()
	if err != nil {
		log.Fatal(err)
	}
	service := NewClusterHealthService()
	r, err := service.HealthCheck("bsdf")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(r)

}
