package client

import (
	"fmt"
	"github.com/KubeOperator/FusionComputeGolangSDK/pkg/client"
	"github.com/KubeOperator/FusionComputeGolangSDK/pkg/cluster"
	"github.com/KubeOperator/FusionComputeGolangSDK/pkg/network"
	"github.com/KubeOperator/FusionComputeGolangSDK/pkg/site"
	"github.com/KubeOperator/FusionComputeGolangSDK/pkg/storage"
	"github.com/KubeOperator/FusionComputeGolangSDK/pkg/vm"
)

func NewFusionComputeClient(vars map[string]interface{}) *fusionComputeClient {
	return &fusionComputeClient{
		Vars: vars,
	}
}

type fusionComputeClient struct {
	Vars map[string]interface{}
}

func (f *fusionComputeClient) ListDatacenter() ([]string, error) {
	c := f.newFusionComputeClient()
	if err := c.Connect(); err != nil {
		return nil, err
	}
	defer func() {
		if err := c.DisConnect(); err != nil {
			fmt.Printf("c.DisConnect()出现了错误：%v\n", err)
		}
	}()
	sm := site.NewManager(c)
	ss, err := sm.ListSite()
	if err != nil {
		return nil, err
	}
	var result []string
	for _, s := range ss {
		result = append(result, s.Name)
	}
	return result, nil
}

func (f *fusionComputeClient) ListClusters() ([]interface{}, error) {
	siteName := f.Vars["datacenter"].(string)
	c := f.newFusionComputeClient()
	if err := c.Connect(); err != nil {
		return nil, err
	}
	defer func() {
		if err := c.DisConnect(); err != nil {
			fmt.Printf("c.DisConnect()出现了错误：%v\n", err)
		}
	}()
	sm := site.NewManager(c)
	ss, err := sm.ListSite()
	if err != nil {
		return nil, err
	}
	siteUri := ""
	for _, s := range ss {
		if s.Name == siteName {
			siteUri = s.Uri
		}
	}
	if siteUri == "" {
		return nil, fmt.Errorf("site %s not found", siteName)
	}
	cm := cluster.NewManager(c, siteUri)
	cs, err := cm.ListCluster()
	if err != nil {
		return nil, err
	}
	var result []interface{}
	for _, cc := range cs {
		ccMeta := make(map[string]interface{})
		ccMeta["cluster"] = cc.Name
		// datastore
		dm := storage.NewManager(c, siteUri)
		ds, err := dm.ListDataStore()
		if err != nil {
			return nil, err
		}
		var dsNames []string
		for _, d := range ds {
			dsNames = append(dsNames, d.Name)
		}
		ccMeta["datastores"] = dsNames
		var templateNames []string
		vmm := vm.NewManager(c, siteUri)
		vms, err := vmm.ListVm(true)
		if err != nil {
			return nil, err
		}
		for _, v := range vms {
			templateNames = append(templateNames, v.Name)
		}
		ccMeta["templates"] = templateNames
		var switchs []map[string]interface{}
		nm := network.NewManager(c, siteUri)
		ss, err := nm.ListDVSwitch()
		if err != nil {
			return nil, err
		}
		for _, s := range ss {
			ssMeta := make(map[string]interface{})
			ssMeta["name"] = s.Name
			ps, err := nm.ListPortGroupBySwitch(s.Uri)
			if err != nil {
				return nil, err
			}
			var portgroups []string
			for _, p := range ps {
				portgroups = append(portgroups, p.Name)
			}
			ssMeta["portgroups"] = portgroups
			switchs = append(switchs, ssMeta)
		}
		ccMeta["switchs"] = switchs
		result = append(result, ccMeta)
	}

	return result, nil
}

func (f *fusionComputeClient) ListTemplates() ([]interface{}, error) {
	return nil, nil
}

func (f *fusionComputeClient) ListFlavors() ([]interface{}, error) {
	return nil, nil
}

func (f *fusionComputeClient) GetIpInUsed(network string) ([]string, error) {
	return []string{}, nil
}

func (f *fusionComputeClient) UploadImage() error {
	return nil
}

func (f *fusionComputeClient) DefaultImageExist() (bool, error) {
	return false, nil
}

func (f *fusionComputeClient) CreateDefaultFolder() error {
	return nil
}

func (f *fusionComputeClient) newFusionComputeClient() client.FusionComputeClient {
	server := f.Vars["server"].(string)
	user := f.Vars["user"].(string)
	password := f.Vars["password"].(string)
	return client.NewFusionComputeClient(server, user, password)
}
