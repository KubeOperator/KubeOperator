package server

import (
	"fmt"
	"github.com/spf13/viper"
	"ko3-gin/pkg/config"
	"ko3-gin/pkg/db"
	"ko3-gin/pkg/logger"
	"ko3-gin/pkg/redis"
	"ko3-gin/pkg/router"
)

type Phase interface {
	Init() error
	PhaseName() string
}

func Phases() []Phase {
	return []Phase{
		//&migrate.InitMigrateDBPhase{
		//	Host:     viper.GetString("db.host"),
		//	Port:     viper.GetInt("db.port"),
		//	Name:     viper.GetString("db.name"),
		//	User:     viper.GetString("db.user"),
		//	Password: viper.GetString("db.password"),
		//},
		&db.InitDBPhase{
			Host:     viper.GetString("db.host"),
			Port:     viper.GetInt("db.port"),
			Name:     viper.GetString("db.name"),
			User:     viper.GetString("db.user"),
			Password: viper.GetString("db.password"),
		},
		&redis.InitRedisPhase{
			Host:       viper.GetString("redis.host"),
			Port:       viper.GetInt("redis.port"),
			DB:         viper.GetInt("db"),
			MaxRetries: 3,
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
	bind := fmt.Sprintf("%s:%d",
		viper.GetString("bind.host"),
		viper.GetInt("bind.port"))
	if err := s.Run(bind); err != nil {
		return err
	}
	return nil
}
