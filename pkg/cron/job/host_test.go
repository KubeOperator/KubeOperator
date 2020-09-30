package job

import (
	"github.com/KubeOperator/KubeOperator/pkg/config"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/spf13/viper"
	"testing"
)

func TestRefreshHostInfo_Run(t *testing.T) {
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
	for i := 0; i < 100; i++ {
		go func() {
			j := NewRefreshHostInfo()
			j.Run()
		}()
	}

	select {}
}
