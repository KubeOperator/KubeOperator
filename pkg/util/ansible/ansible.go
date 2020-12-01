package ansible

import (
	"fmt"
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/util/file"
	uuid "github.com/satori/go.uuid"
	"io"
	"os"
	"path"
)

func CreateAnsibleLogWriter(clusterName string) (string, io.Writer, error) {
	logId := uuid.NewV4().String()
	dirName := path.Join(constant.DefaultAnsibleLogDir, clusterName)
	if !file.Exists(dirName) {
		err := os.MkdirAll(dirName, 0755)
		if err != nil {
			return "", nil, err
		}
	}
	fileName := path.Join(dirName, fmt.Sprintf("%s.log", logId))
	writer, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0755)
	if err != nil {
		return "", nil, err
	}
	return logId, writer, nil
}

func GetAnsibleLogReader(clusterName string, logId string) (io.Reader, error) {
	logPath := path.Join(constant.DefaultAnsibleLogDir, clusterName, fmt.Sprintf("%s.log", logId))
	return os.OpenFile(logPath, os.O_RDONLY, 0755)
}
func GetNodeAnsibleLogReader(clusterName string, nodeName string, logId string) (io.Reader, error) {
	logPath := path.Join(constant.DefaultAnsibleLogDir, clusterName, nodeName, fmt.Sprintf("%s.log", logId))
	return os.OpenFile(logPath, os.O_RDONLY, 0755)
}
