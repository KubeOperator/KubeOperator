package ansible

import (
	"fmt"
	"io"
	"os"
	"path"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/util/file"
	uuid "github.com/satori/go.uuid"
)

func CreateAnsibleLogWriter(clusterName string) (string, string, error) {
	logId := uuid.NewV4().String()
	dirName := path.Join(constant.DefaultAnsibleLogDir, clusterName)
	if !file.Exists(dirName) {
		err := os.MkdirAll(dirName, 0750)
		if err != nil {
			return "", "", err
		}
	}
	fileName := path.Join(dirName, fmt.Sprintf("%s.log", logId))
	return logId, fileName, nil
}

func CreateAnsibleLogWriterWithId(clusterName string, logId string) (string, error) {
	dirName := path.Join(constant.DefaultAnsibleLogDir, clusterName)
	if !file.Exists(dirName) {
		err := os.MkdirAll(dirName, 0750)
		if err != nil {
			return "", err
		}
	}
	fileName := path.Join(dirName, fmt.Sprintf("%s.log", logId))
	return fileName, nil
}

func GetAnsibleLogReader(clusterName string, logId string) (io.Reader, error) {
	logPath := path.Join(constant.DefaultAnsibleLogDir, clusterName, fmt.Sprintf("%s.log", logId))
	return os.OpenFile(logPath, os.O_RDONLY, 0640)
}
func GetNodeAnsibleLogReader(clusterName string, nodeName string, logId string) (io.Reader, error) {
	logPath := path.Join(constant.DefaultAnsibleLogDir, clusterName, nodeName, fmt.Sprintf("%s.log", logId))
	return os.OpenFile(logPath, os.O_RDONLY, 0640)
}
