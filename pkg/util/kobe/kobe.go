package kobe

import (
	"io"

	"github.com/KubeOperator/KubeOperator/api"
	"github.com/spf13/viper"
)

type Interface interface {
	RunPlaybook(name, tag string) (string, error)
	Watch(writer io.Writer, taskId string) error
	GetResult(taskId string) (*api.KobeResult, error)
	SetVar(key string, value string)
}

type Config struct {
	Inventory *api.Inventory
}

type Kobe struct {
	Project   string
	Inventory *api.Inventory
	client    *KobeClient
}

func NewAnsible(c *Config) *Kobe {
	c.Inventory.Vars = map[string]string{}
	host := viper.GetString("kobe.host")
	port := viper.GetInt("kobe.port")
	return &Kobe{
		Project:   "ko",
		Inventory: c.Inventory,
		client:    NewKobeClient(host, port),
	}
}

func (k *Kobe) RunPlaybook(name, tag string) (string, error) {
	result, err := k.client.RunPlaybook(k.Project, name, tag, k.Inventory)
	if err != nil {
		return "", err
	}
	return result.Id, nil
}

func (k *Kobe) SetVar(ke string, va string) {
	k.Inventory.Vars[ke] = va
}

func (k *Kobe) RunAdhoc(pattern, module, param string) (string, error) {
	result, err := k.client.RunAdhoc(pattern, module, param, k.Inventory)
	if err != nil {
		return "", nil
	}
	return result.Id, nil
}

func (k *Kobe) Watch(writer io.Writer, taskId string) error {
	err := k.client.WatchRun(taskId, writer)
	if err != nil {
		return err
	}
	return nil
}

func (k *Kobe) GetResult(taskId string) (*api.KobeResult, error) {
	return k.client.GetResult(taskId)
}
