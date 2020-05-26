package host

import (
	"github.com/KubeOperator/KubeOperator/pkg/config"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	hostModel "github.com/KubeOperator/KubeOperator/pkg/model/host"
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
	item := hostModel.Host{
		Name:         "test",
		Ip:           "172.16.10.63",
		Port:         22,
		CredentialID: "4d73cc9e-2342-4628-b4aa-7461ef04ebca",
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

func TestDelete(t *testing.T) {
	Init()
	err := Delete("test")
	if err != nil {
		t.Fatalf("can not delete item,%s", err)
	}
}
