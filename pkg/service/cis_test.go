package service

import (
	"fmt"
	"github.com/KubeOperator/KubeOperator/pkg/config"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/logger"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/spf13/viper"
	"testing"
)

func TestDo(t *testing.T) {
	config.Init()
	logger.Init()
	dbInit := db.InitDBPhase{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetInt("db.port"),
		Name:     viper.GetString("db.name"),
		User:     viper.GetString("db.user"),
		Password: viper.GetString("db.password"),
	}
	err := dbInit.Init()
	if err != nil {
		fmt.Println(err)
		return
	}
	db.DB.AutoMigrate(model.CisTask{})
	db.DB.AutoMigrate(model.CisResult{})
	service := NewCisService()
	task, err := service.Create("abc")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(task.ID)
	select {}
}
