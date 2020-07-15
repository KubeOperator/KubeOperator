package service

import (
	"github.com/KubeOperator/KubeOperator/pkg/config"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/spf13/viper"
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

func TestClusterIaasService_Init(t *testing.T) {
	Init()
	s := NewClusterIaasService()
	err := s.Init("bbb")
	if err != nil {
		log.Fatal(err)
	}

}
