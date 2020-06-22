package grafana

import (
	"fmt"
	"testing"
)

func GetClient() *Client {
	return NewClient(Config{
		Host:     "localhost",
		Port:     3000,
		Username: "admin",
		Password: "admin",
	})
}

func TestClient_CreateDataSource(t *testing.T) {
	c := GetClient()
	err := c.CreateDataSource(DataSource{
		Name:      "cluster_3",
		Type:      "prometheus",
		Url:       "http://prometheus.apps.zxv.com/",
		Access:    "proxy",
		BasicAuth: false,
	})
	if err != nil {
		fmt.Println(err.Error())
	}
}

func TestClient_CreateDashboard(t *testing.T) {
	c := GetClient()
	err := c.CreateDashboard("cluster_3")
	if err != nil {
		fmt.Println(err.Error())
	}
}

func TestClient_DeleteDataSource(t *testing.T) {
	c := GetClient()
	err := c.DeleteDataSource("cluster_3")
	if err != nil {
		t.Error(err)
	}
}

func TestClient_DeleteDashboard(t *testing.T) {
	c := GetClient()
	err := c.DeleteDashboard("cluster_3")
	if err != nil {
		fmt.Println(err.Error())
		t.Error(err)
	}
}
