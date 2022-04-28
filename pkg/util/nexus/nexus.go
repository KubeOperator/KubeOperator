package nexus

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"time"
)

const (
	BaseUrl = "service/rest"
	RepoUrl = BaseUrl + "/v1/repositories"
)

func CheckConn(username, password, endpoint string) error {
	url := fmt.Sprintf("%s/%s", endpoint, RepoUrl)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	if username != "" && password != "" {
		request.SetBasicAuth(username, password)
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: 2 * time.Second,
	}
	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
	}

	return nil
}
