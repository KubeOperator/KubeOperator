package logger

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	Level   logrus.Level
	Default *logrus.Logger
)

func Init() {
	l := viper.GetString("logging.level")
	level, err := logrus.ParseLevel(l)
	if err != nil {
		Level=logrus.InfoLevel
	} else {
		Level=level
	}
	initDefault()
}

func initDefault()  {
	Default =logrus.New()
	Default.SetLevel(Level)
}

