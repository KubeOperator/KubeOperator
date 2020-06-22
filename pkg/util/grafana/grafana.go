package grafana

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type Interface interface {
	CreateDataSource(source DataSource) error
	DeleteDataSource(name string) error
	CreateDashboard(dataSourceName string) error
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

func NewClient(config Config) *Client {
	return &Client{
		Host:     config.Host,
		Port:     config.Port,
		Username: config.Username,
		Password: config.Password,
	}
}

func (c Client) CreateDataSource(source DataSource) error {
	url := fmt.Sprintf("http://%s:%s@%s:%d/api/datasources/", c.Username, c.Password, c.Host, c.Port)
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

func (c Client) CreateDashboard(dataSourceName string) error {
	dashboard := NewDashboard(dataSourceName)
	req := CreateDashboardRequest{
		Dashboard: *dashboard,
		Overwrite: true,
	}
	url := fmt.Sprintf("http://%s:%s@%s:%d/api/dashboards/db/", c.Username, c.Password, c.Host, c.Port)
	data, err := json.Marshal(&req)
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
