package server

import (
	"fmt"
	"github.com/KubeOperator/KubeOperator/pkg/plugin/xpack"

	"github.com/KubeOperator/KubeOperator/pkg/config"
	"github.com/KubeOperator/KubeOperator/pkg/cron"
	"github.com/KubeOperator/KubeOperator/pkg/data"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/logger"
	"github.com/KubeOperator/KubeOperator/pkg/migrate"
	"github.com/KubeOperator/KubeOperator/pkg/plugin"
	"github.com/KubeOperator/KubeOperator/pkg/router"
	"github.com/kataras/iris/v12"
	"github.com/spf13/viper"
)

type Phase interface {
	Init() error
	PhaseName() string
}

func Phases() []Phase {
	return []Phase{
		&db.InitDBPhase{
			Host:     viper.GetString("db.host"),
			Port:     viper.GetInt("db.port"),
			Name:     viper.GetString("db.name"),
			User:     viper.GetString("db.user"),
			Password: viper.GetString("db.password"),
		},
		&migrate.InitMigrateDBPhase{
			Host:     viper.GetString("db.host"),
			Port:     viper.GetInt("db.port"),
			Name:     viper.GetString("db.name"),
			User:     viper.GetString("db.user"),
			Password: viper.GetString("db.password"),
		},
		&data.InitDataPhase{},
		&plugin.InitPluginDBPhase{},
		&cron.InitCronPhase{
			Enable: viper.GetBool("cron.enable"),
		},
	}
}

func Start() error {
	config.Init()
	logger.Init()
	var log = logger.Default
	phases := Phases()
	for _, phase := range phases {
		if err := phase.Init(); err != nil {
			log.Errorf("start phase [%s] failed reason: %s",
				phase.PhaseName(), err.Error())
			return err
		}
		log.Infof("start phase [%s] success", phase.PhaseName())
	}
	s := router.Server()
	// load xpack plugin must behead router init,so can not create an phase for it.
	if err := xpack.LoadXpackPlugin(); err != nil {
		log.Error("xpack load failed, xpack can not be registered")
	}
	bind := fmt.Sprintf("%s:%d",
		viper.GetString("bind.host"),
		viper.GetInt("bind.port"))

	if err := s.Run(iris.Addr(bind)); err != nil {
		return err
	}
	return nil
}
