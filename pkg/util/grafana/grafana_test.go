package grafana

import (
	"fmt"
	"github.com/KubeOperator/KubeOperator/pkg/config"
	"testing"
)

func GetClient() *Client {
	config.Init()
	return NewClient()
}

func TestClient_CreateDashboard(t *testing.T) {
	c := GetClient()
	url, err := c.CreateDashboard("zxv")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(url)
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
	err := c.DeleteDashboard("zxv")
	if err != nil {
		fmt.Println(err.Error())
		t.Error(err)
	}
}
