package client

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"

	"github.com/Azure/azure-storage-blob-go/azblob"
)

type azureClient struct {
	Vars       map[string]interface{}
	ServiceURL azblob.ServiceURL
}

func NewAzureClient(vars map[string]interface{}) (*azureClient, error) {
	var accountName string
	var accountKey string
	var endpoint string
	if _, ok := vars["accountName"]; ok {
		accountName = vars["accountName"].(string)
	} else {
		return nil, errors.New(ParamEmpty)
	}
	if _, ok := vars["accountKey"]; ok {
		accountKey = vars["accountKey"].(string)
	} else {
		return nil, errors.New(ParamEmpty)
	}
	if _, ok := vars["endpoint"]; ok {
		endpoint = vars["endpoint"].(string)
	} else {
		return nil, errors.New(ParamEmpty)
	}
	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		return nil, err
	}
	p := azblob.NewPipeline(credential, azblob.PipelineOptions{})
	u, err := url.Parse(fmt.Sprintf("https://%s."+endpoint, accountName))
	if err != nil {
		return nil, err
	}
	serviceURL := azblob.NewServiceURL(*u, p)
	return &azureClient{
		Vars:       vars,
		ServiceURL: serviceURL,
	}, nil
}

func (azure azureClient) ListBuckets() ([]interface{}, error) {

	response, err := azure.ServiceURL.ListContainersSegment(context.Background(), azblob.Marker{}, azblob.ListContainersSegmentOptions{})
	if err != nil {
		return nil, err
	}
	var result []interface{}
	for _, bucket := range response.ContainerItems {
		result = append(result, bucket.Name)
	}
	return result, nil
}

func (azure azureClient) Exist(path string) (bool, error) {
	return true, nil
}

func (azure azureClient) Delete(path string) (bool, error) {
	containerURL, err := azure.getBucket()
	if err != nil {
		return false, err
	}
	blobURL := containerURL.NewBlockBlobURL(path)
	_, err = blobURL.Delete(context.Background(), azblob.DeleteSnapshotsOptionInclude, azblob.BlobAccessConditions{})
	if err != nil {
		if serr, ok := err.(azblob.StorageError); ok {
			if serr.Response().StatusCode == 404 {
				return true, nil
			}
		} else {
			return false, serr
		}
		return false, err
	}
	return true, nil
}

func (azure azureClient) Upload(src, target string) (bool, error) {
	containerURL, err := azure.getBucket()
	if err != nil {
		return false, err
	}
	blobURL := containerURL.NewBlockBlobURL(target)
	file, err := os.Open(src)
	if err != nil {
		return false, err
	}
	_, err = azblob.UploadFileToBlockBlob(context.Background(), file, blobURL, azblob.UploadToBlockBlobOptions{
		Parallelism: 16,
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (azure azureClient) Download(src, target string) (bool, error) {
	containerURL, err := azure.getBucket()
	if err != nil {
		return false, err
	}
	blobURL := containerURL.NewBlockBlobURL(src)
	downloadResponse, err := blobURL.Download(context.Background(), 0, azblob.CountToEnd, azblob.BlobAccessConditions{}, false)
	if err != nil {
		return false, err
	}
	bodyStream := downloadResponse.Body(azblob.RetryReaderOptions{MaxRetryRequests: 20})
	downloadedData := bytes.Buffer{}
	_, err = downloadedData.ReadFrom(bodyStream)
	bodyStream.Close()
	if err != nil {
		return false, err
	}
	err = ioutil.WriteFile(target, downloadedData.Bytes(), 0750)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (azure *azureClient) getBucket() (*azblob.ContainerURL, error) {
	if _, ok := azure.Vars["bucket"]; ok {
		containerURL := azure.ServiceURL.NewContainerURL(azure.Vars["bucket"].(string))
		return &containerURL, nil
	} else {
		return nil, errors.New(ParamEmpty)
	}
}
