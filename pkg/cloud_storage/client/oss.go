package client

import (
	"errors"

	"github.com/KubeOperator/KubeOperator/pkg/util/escape"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

var (
	ParamEmpty = "PARAM_EMPTY"
)

type ossClient struct {
	Vars   map[string]interface{}
	client oss.Client
}

func NewOssClient(vars map[string]interface{}) (*ossClient, error) {
	var endpoint string
	var accessKey string
	var secretKey []byte
	if _, ok := vars["endpoint"]; ok {
		endpoint, ok = vars["endpoint"].(string)
		if !ok {
			return nil, errors.New("type aassertion failed")
		}
	} else {
		return nil, errors.New(ParamEmpty)
	}
	if _, ok := vars["accessKey"]; ok {
		accessKey, ok = vars["accessKey"].(string)
		if !ok {
			return nil, errors.New("type aassertion failed")
		}
	} else {
		return nil, errors.New(ParamEmpty)
	}
	if _, ok := vars["secretKey"]; ok {
		secretKey = escape.GetByte(vars["secretKey"])
	} else {
		return nil, errors.New(ParamEmpty)
	}
	client, err := oss.New(endpoint, accessKey, string(secretKey))
	if err != nil {
		return nil, err
	}
	escape.Clean(string(secretKey))

	return &ossClient{
		Vars:   vars,
		client: *client,
	}, nil
}

func (oss ossClient) ListBuckets() ([]interface{}, error) {
	response, err := oss.client.ListBuckets()
	if err != nil {
		return nil, err
	}
	var result []interface{}
	for _, bucket := range response.Buckets {
		result = append(result, bucket.Name)
	}
	return result, err
}

func (oss ossClient) Exist(path string) (bool, error) {
	bucket, err := oss.GetBucket()
	if err != nil {
		return false, err
	}
	return bucket.IsObjectExist(path)

}

func (oss ossClient) Delete(path string) (bool, error) {
	bucket, err := oss.GetBucket()
	if err != nil {
		return false, err
	}
	err = bucket.DeleteObject(path)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (oss ossClient) Upload(src, target string) (bool, error) {
	bucket, err := oss.GetBucket()
	if err != nil {
		return false, err
	}
	err = bucket.PutObjectFromFile(target, src)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (oss ossClient) Download(src, target string) (bool, error) {
	bucket, err := oss.GetBucket()
	if err != nil {
		return false, err
	}
	err = bucket.GetObjectToFile(src, target)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (oss *ossClient) GetBucket() (*oss.Bucket, error) {
	if _, ok := oss.Vars["bucket"]; ok {
		bucket, err := oss.client.Bucket(oss.Vars["bucket"].(string))
		if err != nil {
			return nil, err
		}
		return bucket, nil
	} else {
		return nil, errors.New(ParamEmpty)
	}
}
