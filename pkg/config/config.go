package config

import (
	"github.com/spf13/viper"
)

func Init() {
	viper.SetConfigName("app")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/ko")
	_ = viper.ReadInConfig()
}
