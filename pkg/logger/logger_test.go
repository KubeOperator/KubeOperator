package logger

import (
	"io"
	"os"
	"testing"
	"time"

	"github.com/KubeOperator/KubeOperator/pkg/config"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func TestInit(t *testing.T) {
	config.Init()
	log := logrus.New()
	path := "/tmp/mylog/ko_log"
	l := viper.GetString("logging.level")
	outPut := viper.GetString("logging.out_put")
	maxAge := viper.GetInt("logging.max_age")
	rotationTime := viper.GetInt("logging.rotation")
	level, err := logrus.ParseLevel(l)
	if err != nil && l == "" {
		log.SetLevel(logrus.InfoLevel)
	}
	log.SetLevel(level)
	log.SetReportCaller(true)
	log.SetFormatter(new(MineFormatter))
	writer, _ := rotatelogs.New(
		path+".%Y%m%d%H%M",
		rotatelogs.WithLinkName(path),
		rotatelogs.WithMaxAge(time.Duration(maxAge)*time.Second),
		rotatelogs.WithRotationTime(time.Duration(rotationTime)*time.Second),
	)
	switch outPut {
	case "file":
	case "fileAndStd":
		writers := []io.Writer{writer, os.Stdout}
		fileAndStdoutWriter := io.MultiWriter(writers...)
		log.SetOutput(fileAndStdoutWriter)
	case "std":
		log.SetOutput(writer)
	default:
		writers := []io.Writer{writer, os.Stdout}
		fileAndStdoutWriter := io.MultiWriter(writers...)
		log.SetOutput(fileAndStdoutWriter)
	}

	for {
		log.WithFields(logrus.Fields{
			"animal": "walrus",
			"size":   10,
		}).Info("A group of walrus emerges from the ocean")
		time.Sleep(time.Duration(2) * time.Second)
	}
}
