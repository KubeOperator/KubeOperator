package logger

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var Log *logrus.Logger

type MineFormatter struct{}

const TimeFormat = "2006-01-02 15:04:05"

func (s *MineFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var cstSh, _ = time.LoadLocation("Asia/Shanghai")
	msg := fmt.Sprintf("[%s] [%s] %s (%s: %d) {%v} \n", time.Now().In(cstSh).Format(TimeFormat), strings.ToUpper(entry.Level.String()), entry.Message, entry.Caller.Function, entry.Caller.Line, entry.Data)
	return []byte(msg), nil
}

func Init() {
	log := logrus.New()
	path := "/var/ko/data/logs/log"

	l := viper.GetString("logging.level")
	outPut := viper.GetString("logging.out_put")
	maxAge := viper.GetInt("logging.max_age")
	rotationTime := viper.GetInt("logging.rotation")

	level, err := logrus.ParseLevel(l)
	if err != nil && l == "" {
		log.SetLevel(logrus.InfoLevel)
	} else {
		log.SetLevel(level)
	}
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
	Log = log
}
