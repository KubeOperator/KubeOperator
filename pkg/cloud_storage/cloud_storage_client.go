package cloud_storage

import (
	"errors"
	"github.com/KubeOperator/KubeOperator/pkg/cloud_storage/client"
	"github.com/KubeOperator/KubeOperator/pkg/constant"
)

var (
	NotSupport = "NOT_SUPPORT"
)

type CloudStorageClient interface {
	ListBuckets() ([]interface{}, error)
	Exist(path string) (bool, error)
	Delete(path string) (bool, error)
	Upload(src, target string) (bool, error)
	Download(src, target string) (bool, error)
}

func NewCloudStorageClient(vars map[string]string) (CloudStorageClient, error) {
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
