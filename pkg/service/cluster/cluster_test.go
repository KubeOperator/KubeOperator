package cluster

import (
	"github.com/spf13/viper"
	"ko3-gin/pkg/config"
	"ko3-gin/pkg/db"
	clusterModel "ko3-gin/pkg/model/cluster"
	"ko3-gin/pkg/model/common"
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

func TestSave(t *testing.T) {
	Init()
	item := clusterModel.Cluster{
		BaseModel: common.BaseModel{
			Name: "test",
		},
	}
	err := db.DB.Create(&item).Error
	if err != nil {
		log.Fatalf("can not create item,%s", err)
	}
}
