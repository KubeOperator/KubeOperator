package ansible

import (
	"fmt"
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/util/file"
	"io"
	"os"
	"path"
)

func CreateAnsibleLogWriter(clusterName string, logId string) (io.Writer, error) {
	dirName := path.Join(constant.DefaultAnsibleLogDir, clusterName)
	if !file.Exists(dirName) {
		_ = os.MkdirAll(dirName, 0755)
	}
	fileName := path.Join(dirName, fmt.Sprintf("%s.log", logId))
	return os.OpenFile(fileName, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0755)
}

func GetAnsibleLogReader(clusterName string, logId string) (io.Reader, error) {
	logPath := path.Join(constant.DefaultAnsibleLogDir, clusterName, fmt.Sprintf("%s.log", logId))
	return os.OpenFile(logPath, os.O_RDONLY, 0755)
}
