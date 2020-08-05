package backup_acoount

import (
	"errors"
	"github.com/KubeOperator/KubeOperator/pkg/backup_acoount/client"
	"github.com/KubeOperator/KubeOperator/pkg/constant"
)

var (
	NotSupport = "NOT_SUPPORT"
)

type BackupAccountClient interface {
	ListBuckets() ([]interface{}, error)
	Exist(path string) (bool, error)
	Delete(path string) (bool, error)
	Upload(src, target string) (bool, error)
	Download(src, target string) (bool, error)
}

func NewBackupAccountClient(vars map[string]string) (BackupAccountClient, error) {
	if vars["type"] == constant.Azure {
		return client.NewAzureClient(vars)
	}
	if vars["type"] == constant.S3 {
		return client.NewS3Client(vars)
	}
	if vars["type"] == constant.OSS {
		return client.NewOssClient(vars)
	}
	return nil, errors.New(NotSupport)
}
