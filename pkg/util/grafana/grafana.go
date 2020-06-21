package grafana

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Interface interface {
	CreateDataSource(source DataSource) error
	CreateDashboard(dashboard Dashboard) error
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
	var buffer bytes.Buffer
	_, err = buffer.Read(data)
	if err != nil && err != io.EOF {
		return err
	}
	resp, err := http.Post(url, "application/json", &buffer)
	if err != nil {
		return err
	}
	fmt.Println(resp.StatusCode)
	return nil
}
