package cluster

import (
	"fmt"
	"github.com/KubeOperator/KubeOperator/pkg/config"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	clusterModel "github.com/KubeOperator/KubeOperator/pkg/model/cluster"
	"github.com/spf13/viper"
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
	for i := 0; i < 100; i++ {
		item := clusterModel.Cluster{
			Name: "test" + string(i),
			Spec: clusterModel.Spec{
				Version:     "v1.18.2",
				NetworkType: "calico",
				ClusterCIDR: "172.16.10.142/8",
				ServiceCIDR: "172.16.10.142/8",
			},
		}
		err := Save(&item)
		if err != nil {
			t.Fatalf("can not create item,%s", err)
		}
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

func TestDelete(t *testing.T) {
	Init()
	err := Delete("test")
	if err != nil {
		t.Fatalf("can not delete item,%s", err)
	}
}
