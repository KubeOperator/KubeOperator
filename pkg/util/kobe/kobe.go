package kobe

import (
	"github.com/KubeOperator/kobe/api"
	kobeClient "github.com/KubeOperator/kobe/pkg/client"
	"io"
)

type Interface interface {
	RunPlaybook(name string) (string, error)
	Watch(writer io.Writer, taskId string) error
	GetResult(taskId string) (*api.Result, error)
	SetVar(key string, value string)
}

type Config struct {
	Host      string
	Port      int
	Inventory api.Inventory
}

type Kobe struct {
	Project   string
	Inventory api.Inventory
	client    *kobeClient.KobeClient
}

func NewAnsible(c *Config) *Kobe {
	c.Inventory.Vars = map[string]string{}
	return &Kobe{
		Project:   "ko",
		Inventory: c.Inventory,
		client:    kobeClient.NewKobeClient(c.Host, c.Port),
	}
}

func (k *Kobe) RunPlaybook(name string) (string, error) {
	result, err := k.client.RunPlaybook(k.Project, name, k.Inventory)
	if err != nil {
		return "", err
	}
	return result.Id, nil
}

func (k *Kobe) SetVar(key string, value string) {
	k.Inventory.Vars[key] = value
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

func (a *Kobe) GetResult(taskId string) (*api.Result, error) {
	return a.client.GetResult(taskId)
}
