package kobe

import (
	"github.com/KubeOperator/kobe/api"
	"os"
	"testing"
)

func TestKobe_RunAdhoc(t *testing.T) {
	ansible := NewAnsible(&Config{
		Host: "localhost",
		Port: 8081,
		Inventory: api.Inventory{
			Hosts: []*api.Host{
				{
					Ip:       "47.244.14.164",
					Name:     "test",
					Port:     22,
					User:     "root",
					Password: "Calong@2015",
					Vars:     map[string]string{},
				},
			},
			Groups: []*api.Group{
				{
					Name:     "master",
					Children: []string{},
					Vars:     map[string]string{},
					Hosts:    []string{"test"},
				},
			},
		},
	})
	resultId, err := ansible.RunAdhoc("master", "setup", "")
	if err != nil {
		t.Fatal(err)
	}
	err = ansible.Watch(os.Stdout, resultId)
	if err != nil {
		t.Fatal(err)
	}
	res, err := ansible.GetResult(resultId)
	if err != nil {
		t.Fatal(err)
	}
	println(res.Content)
}

func TestKobe_RunPlaybook(t *testing.T) {
	ansible := NewAnsible(&Config{
		Host: "localhost",
		Port: 8081,
		Inventory: api.Inventory{
			Hosts: []*api.Host{
				{
					Ip:       "47.244.14.164",
					Name:     "test",
					Port:     22,
					User:     "root",
					Password: "Calong@2015",
					Vars:     map[string]string{},
				},
			},
			Groups: []*api.Group{
				{
					Name:     "master",
					Children: []string{},
					Vars:     map[string]string{},
					Hosts:    []string{"test"},
				},
			},
		},
	})
	resultId, err := ansible.RunPlaybook("test.yml")
	if err != nil {
		t.Fatal(err)
	}
	err = ansible.Watch(os.Stdout, resultId)
	if err != nil {
		t.Fatal(err)
	}
	res, err := ansible.GetResult(resultId)
	println(res.Content)
}
