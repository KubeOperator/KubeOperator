package job

import (
	"testing"

	"github.com/KubeOperator/KubeOperator/pkg/config"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/logger"
	"github.com/spf13/viper"
)

func TestClusterEvent(t *testing.T) {
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
		logger.Log.Fatal(err)
	}
	j := NewClusterEvent()
	j.Run()
	select {}
}
