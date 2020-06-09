package logger

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	Default = logrus.New()
)

func Init() {
	l := viper.GetString("logging.level")
	level, err := logrus.ParseLevel(l)
	if err != nil && l == "" {
		Default.SetLevel(logrus.InfoLevel)
	}
	logrus.SetLevel(level)
}
