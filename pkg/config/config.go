package config

import (
	"os"
	"strings"

	"github.com/spf13/viper"
)

func Init() {
	viper.SetConfigName("app")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/ko")
	_ = viper.ReadInConfig()
	splitOsEnv()
}

func splitOsEnv() {
	for i := range os.Environ() {
		ks := strings.Split(os.Environ()[i], "=")
		if len(ks) < 2 {
			continue
		}
		key := ks[0]
		value := ks[1]
		if strings.HasPrefix(key, "KO_") {
			cfk := strings.Replace(strings.ToLower(strings.Replace(key, "KO_", "", -1)), "_", ".", -1)
			viper.Set(cfk, value)
		}
	}
}
