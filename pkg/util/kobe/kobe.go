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
	return &Kobe{
		Project:   "ko",
		Inventory: c.Inventory,
		client:    kobeClient.NewKobeClient(c.Host, c.Port),
	}
}

func (a *Kobe) RunPlaybook(name string) (string, error) {
	result, err := a.client.RunPlaybook(a.Project, name, a.Inventory)
	if err != nil {
		return "", err
	}
	return result.Id, nil
}

func (a *Kobe) Watch(writer io.Writer, taskId string) error {
	err := a.client.WatchRunPlaybook(taskId, writer)
	if err != nil {
		return err
	}
	return nil
}

func (a *Kobe) GetResult(taskId string) (*api.Result, error) {
	return a.client.GetResult(taskId)
}
