package kobe

import (
	"testing"
)

func TestKobe_RunAdhoc(t *testing.T) {
	//ansible := NewAnsible(&Config{
	//	Host: "localhost",
	//	Port: 8080,
	//	Inventory: api.Inventory{
	//		Hosts: []*api.Host{
	//			{
	//				Ip:       "47.244.14.164",
	//				Name:     "test",
	//				Port:     22,
	//				User:     "",
	//				Password: "",
	//				Vars:     map[string]string{},
	//			},
	//		},
	//		Groups: []*api.Group{
	//			{
	//				Name:     "master",
	//				Children: []string{},
	//				Vars:     map[string]string{},
	//				Hosts:    []string{"test"},
	//			},
	//		},
	//	},
	//})
	//resultId, err := ansible.RunAdhoc("master", "setup", "")
	//if err != nil {
	//	t.Fatal(err)
	//}
	//err = ansible.Watch(os.Stdout, resultId)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//res, err := ansible.GetResult(resultId)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//println(res.Content)
}

//func TestKobe_RunPlaybook(t *testing.T) {
//	ansible := NewAnsible(&Config{
//		Host: "localhost",
//		Port: 8081,
//		Inventory: api.Inventory{
//			Hosts: []*api.Host{
//				{
//					Ip:       "172.16.10.142",
//					Name:     "test",
//					Port:     22,
//					User:     "",
//					Password: "",
//					Vars:     map[string]string{},
//				},
//			},
//			Groups: []*api.Group{
//				{
//					Name:     "master",
//					Children: []string{},
//					Vars:     map[string]string{},
//					Hosts:    []string{"test"},
//				},
//			},
//		},
//	})
//
//	//ansible.SetVar("bin_dir","/usr/local/bin")
//
//	resultId, err := ansible.RunPlaybook("01-base.yml")
//	if err != nil {
//		t.Fatal(err)
//	}
//	err = ansible.Watch(os.Stdout, resultId)
//	if err != nil {
//		t.Fatal(err)
//	}
//	res, err := ansible.GetResult(resultId)
//	if err != nil {
//		t.Fatal(err)
//	}
//	f, _ := os.Create("a.json")
//	_, _ = f.Write([]byte(res.Content))
//	result, err := ParseResult(res.Content)
//	if err != nil {
//		t.Fatal(err)
//	}
//	result.GatherFailedInfo()
//	f.Close()
//}
