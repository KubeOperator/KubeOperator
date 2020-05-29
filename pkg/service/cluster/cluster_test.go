package cluster

import (
	"fmt"
	"github.com/KubeOperator/KubeOperator/pkg/config"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	clusterModel "github.com/KubeOperator/KubeOperator/pkg/model/cluster"
	hostModel "github.com/KubeOperator/KubeOperator/pkg/model/host"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/viper"
	"log"
	"testing"
)

func Init() {
	config.Init()
	phase := db.InitDBPhase{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetInt("db.port"),
		Name:     viper.GetString("db.name"),
		User:     viper.GetString("db.user"),
		Password: viper.GetString("db.password"),
	}
	err := phase.Init()
	if err != nil {
		log.Fatalf("can not init db,%s", err)
	}
}

func TestSave(t *testing.T) {
	Init()
	item := clusterModel.Cluster{
		Name: "test",
		Spec: clusterModel.Spec{
			Version:     "v1.18.2",
			NetworkType: "calico",
			ClusterCIDR: "172.16.10.142/8",
			ServiceCIDR: "172.16.10.142/8",
		},
		Nodes: []clusterModel.Node{
			{
				ID:   uuid.NewV4().String(),
				Name: "node-1",
				Role: "master",
				Host: hostModel.Host{Name: "test"},
			},
		},
	}
	err := Save(&item)
	if err != nil {
		t.Fatalf("can not create item,%s", err)
	}
}

func TestList(t *testing.T) {
	Init()
	items, err := List()
	if err != nil {
		t.Fatalf("can not list item,%s", err)
	}
	t.Log(items)
}

func TestPage(t *testing.T) {
	Init()
	items, total, err := Page(1, 10)
	if err != nil {
		t.Fatalf("can not page item,%s", err)
	}
	t.Log(items)
	t.Log(total)
}

func TestGet(t *testing.T) {
	Init()
	c, err := Get("test")
	if err != nil {
		t.Fatalf("get item error: %s", err.Error())
	}
	fmt.Println(c.Spec)
}

func TestGetStatus(t *testing.T) {
	Init()
	c, err := GetStatus("test")
	if err != nil {
		t.Fatalf("get item error: %s", err.Error())
	}
	fmt.Println(c)
}

func TestDelete(t *testing.T) {
	Init()
	err := Delete("test")
	if err != nil {
		t.Fatalf("can not delete item,%s", err)
	}
}

type Cluster struct {
	ID string
}

func TestInitCluster(t *testing.T) {
	Init()
	inventory, err := GetKobeInventory("test")
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Println(inventory)

}
