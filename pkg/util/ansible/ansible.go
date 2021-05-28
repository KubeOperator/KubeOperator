package ansible

import (
	"fmt"
	"io"
	"os"
	"path"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/util/file"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

func CreateAnsibleLogWriter(clusterName string) (string, io.Writer, error) {
	logId := uuid.NewV4().String()
	dirName := path.Join(constant.DefaultAnsibleLogDir, clusterName)
	if !file.Exists(dirName) {
		err := os.MkdirAll(dirName, 0755)
		if err != nil {
			return "", nil, errors.Wrap(err, fmt.Sprintf("create ansible log file failed: %v", err))
		}
	}
	fileName := path.Join(dirName, fmt.Sprintf("%s.log", logId))
	writer, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0755)
	if err != nil {
		return "", nil, errors.Wrap(err, fmt.Sprintf("open ansible log file failed: %v", err))
	}
	return logId, writer, nil
}

func CreateAnsibleLogWriterWithId(clusterName string, logId string) (io.Writer, error) {
	dirName := path.Join(constant.DefaultAnsibleLogDir, clusterName)
	if !file.Exists(dirName) {
		err := os.MkdirAll(dirName, 0755)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("open ansible log file with id failed: %v", err))
		}
	}
	fileName := path.Join(dirName, fmt.Sprintf("%s.log", logId))
	writer, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0755)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("open ansible log file with id failed: %v", err))
	}
	return writer, nil
}

func GetAnsibleLogReader(clusterName string, logId string) (io.Reader, error) {
	logPath := path.Join(constant.DefaultAnsibleLogDir, clusterName, fmt.Sprintf("%s.log", logId))
	result, err := os.OpenFile(logPath, os.O_RDONLY, 0755)
	if err != nil {
		return result, errors.Wrap(err, fmt.Sprintf("get ansible log of cluster: %s logid: %s failed: %v", clusterName, logId, err))
	}
	return result, nil
}
func GetNodeAnsibleLogReader(clusterName string, nodeName string, logId string) (io.Reader, error) {
	logPath := path.Join(constant.DefaultAnsibleLogDir, clusterName, nodeName, fmt.Sprintf("%s.log", logId))
	result, err := os.OpenFile(logPath, os.O_RDONLY, 0755)
	if err != nil {
		return result, errors.Wrap(err, fmt.Sprintf("get ansible log of cluster: %s node: %s logid: %s failed: %v", clusterName, nodeName, logId, err))
	}
	return result, nil
}
