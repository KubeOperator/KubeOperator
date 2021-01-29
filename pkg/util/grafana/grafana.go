package grafana

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
	"strings"
)

type Interface interface {
	CreateDataSource(name string, url string) error
	DeleteDataSource(name string) error
	CreateDashboard(dataSourceName string) (string, error)
	DeleteDashboard(name string) error
}

type Config struct {
	Host     string
	Port     int
	Username string
	Password string
}

type Client struct {
	Host     string
	Port     int
	Username string
	Password string
}

func NewClient() *Client {
	return &Client{
		Host:     viper.GetString("grafana.host"),
		Port:     viper.GetInt("grafana.port"),
		Username: "admin",
		Password: "admin",
	}
}

func (c Client) CreateDataSource(name string, prometheusUrl string) error {
	url := fmt.Sprintf("http://%s:%s@%s:%d/api/datasources/", c.Username, c.Password, c.Host, c.Port)
	source := NewDataSource(name, prometheusUrl)
	data, err := json.Marshal(&source)
	if err != nil {
		return err
	}
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		var buf []byte
		_, _ = resp.Body.Read(buf)
		return errors.New(string(buf))
	}
	return nil
}

func (c Client) DeleteDataSource(name string) error {
	url := fmt.Sprintf("http://%s:%s@%s:%d/api/datasources/name/%s", c.Username, c.Password, c.Host, c.Port, name)
	httpClient := http.Client{}
	req, err := http.NewRequest(http.MethodDelete, url, strings.NewReader(""))
	if err != nil {
		return err
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		var buf []byte
		_, _ = resp.Body.Read(buf)
		return errors.New(string(buf))
	}
	return nil
}

func (c Client) CreateDashboard(dataSourceName string) (string, error) {
	dashboard := NewDashboard(dataSourceName)
	req := CreateDashboardRequest{
		Dashboard: *dashboard,
		Overwrite: true,
	}
	url := fmt.Sprintf("http://%s:%s@%s:%d/api/dashboards/db/", c.Username, c.Password, c.Host, c.Port)
	data, err := json.Marshal(&req)
	if err != nil {
		return "", err
	}
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	msg := string(body)
	if resp.StatusCode != 200 {
		return "", errors.New(msg)
	}
	var respMap map[string]interface{}
	err = json.Unmarshal([]byte(msg), &respMap)
	if err != nil {
		return "", err
	}
	u := respMap["url"].(string)
	return u, nil
}

func (c Client) DeleteDashboard(name string) error {
	url := fmt.Sprintf("http://%s:%s@%s:%d/api/dashboards/db/%s", c.Username, c.Password, c.Host, c.Port, name)
	httpClient := http.Client{}
	req, err := http.NewRequest(http.MethodDelete, url, strings.NewReader(""))
	if err != nil {
		return err
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		var buf []byte
		_, _ = resp.Body.Read(buf)
		return errors.New(string(buf))
	}
	return nil
}
