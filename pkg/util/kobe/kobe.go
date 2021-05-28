package kobe

import (
	"fmt"
	"io"

	"github.com/KubeOperator/kobe/api"
	kobeClient "github.com/KubeOperator/kobe/pkg/client"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type Interface interface {
	RunPlaybook(name, tag string) (string, error)
	Watch(writer io.Writer, taskId string) error
	GetResult(taskId string) (*api.Result, error)
	SetVar(key string, value string)
}

type Config struct {
	Inventory *api.Inventory
}

type Kobe struct {
	Project   string
	Inventory *api.Inventory
	client    *kobeClient.KobeClient
}

func NewAnsible(c *Config) *Kobe {
	c.Inventory.Vars = map[string]string{}
	host := viper.GetString("kobe.host")
	port := viper.GetInt("kobe.port")
	return &Kobe{
		Project:   "ko",
		Inventory: c.Inventory,
		client:    kobeClient.NewKobeClient(host, port),
	}
}

func (k *Kobe) RunPlaybook(name, tag string) (string, error) {
	result, err := k.client.RunPlaybook(k.Project, name, tag, k.Inventory)
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("ansible run playbook failed: %v", err))
	}
	return result.Id, nil
}

func (k *Kobe) SetVar(key string, value string) {
	k.Inventory.Vars[key] = value
}

func (k *Kobe) RunAdhoc(pattern, module, param string) (string, error) {
	result, err := k.client.RunAdhoc(pattern, module, param, k.Inventory)
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("ansible run adhoc failed: %v", err))
	}
	return result.Id, nil
}

func (k *Kobe) Watch(writer io.Writer, taskId string) error {
	if err := k.client.WatchRun(taskId, writer); err != nil {
		return errors.Wrap(err, fmt.Sprintf("ansible run watch failed: %v", err))
	}
	return nil
}

func (k *Kobe) GetResult(taskId string) (*api.Result, error) {
	result, err := k.client.GetResult(taskId)
	if err != nil {
		return result, errors.Wrap(err, fmt.Sprintf("ansible get result failed: %v", err))
	}
	return result, nil
}
