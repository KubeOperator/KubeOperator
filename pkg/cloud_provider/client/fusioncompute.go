package client

import (
	"errors"
	"fmt"
	"time"

	"github.com/KubeOperator/FusionComputeGolangSDK/pkg/client"
	"github.com/KubeOperator/FusionComputeGolangSDK/pkg/cluster"
	"github.com/KubeOperator/FusionComputeGolangSDK/pkg/network"
	"github.com/KubeOperator/FusionComputeGolangSDK/pkg/site"
	"github.com/KubeOperator/FusionComputeGolangSDK/pkg/storage"
	"github.com/KubeOperator/FusionComputeGolangSDK/pkg/task"
	"github.com/KubeOperator/FusionComputeGolangSDK/pkg/vm"
	"github.com/KubeOperator/KubeOperator/pkg/constant"
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
			log.Errorf("fusionComputeClient DisConnect failed, error: %s", err.Error())
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
			log.Errorf("fusionComputeClient DisConnect failed, error: %s", err.Error())
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
	siteName := f.Vars["datacenter"].(string)
	var result []string
	c := f.newFusionComputeClient()
	if err := c.Connect(); err != nil {
		return nil, err
	}
	defer func() {
		if err := c.DisConnect(); err != nil {
			log.Errorf("fusionComputeClient DisConnect failed, error: %s", err.Error())
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
	vmm := vm.NewManager(c, siteUri)
	vms, err := vmm.ListVm(false)
	if err != nil {
		return nil, err
	}
	for _, v := range vms {
		result = append(result, v.VmConfig.Nics[0].Ip)
	}
	return result, nil
}

func (f *fusionComputeClient) UploadImage() error {
	siteName, ok := f.Vars["datacenter"].(string)
	if !ok {
		return errors.New("type aassertion failed")
	}
	clusterName, ok := f.Vars["cluster"].(string)
	if !ok {
		return errors.New("type aassertion failed")
	}
	datastoreName, ok := f.Vars["datastore"].(string)
	if !ok {
		return errors.New("type aassertion failed")
	}
	portgroupName, ok := f.Vars["portgroup"].(string)
	if !ok {
		return errors.New("type aassertion failed")
	}
	ovfPath, ok := f.Vars["ovfPath"].(string)
	if !ok {
		return errors.New("type aassertion failed")
	}
	c := f.newFusionComputeClient()
	if err := c.Connect(); err != nil {
		return err
	}
	defer func() {
		if err := c.DisConnect(); err != nil {
			log.Errorf("fusionComputeClient DisConnect failed, error: %s", err.Error())
		}
	}()
	sm := site.NewManager(c)
	ss, err := sm.ListSite()
	if err != nil {
		return err
	}
	siteUri := ""
	for _, s := range ss {
		if s.Name == siteName {
			siteUri = s.Uri
		}
	}

	cm := cluster.NewManager(c, siteUri)
	cs, err := cm.ListCluster()
	if err != nil {
		return err
	}
	clusterUrn := ""
	for _, cc := range cs {
		if cc.Name == clusterName {
			clusterUrn = cc.Urn
		}
	}
	dm := storage.NewManager(c, siteUri)
	ds, err := dm.ListDataStore()
	if err != nil {
		return err
	}
	datastoreUrn := ""
	for _, d := range ds {
		if d.Name == datastoreName {
			datastoreUrn = d.Urn
		}
	}
	nm := network.NewManager(c, siteUri)
	dss, err := nm.ListDVSwitch()
	if err != nil {
		return err
	}
	portGroupUrn := ""
	for _, ds := range dss {
		ps, err := nm.ListPortGroupBySwitch(ds.Uri)
		if err != nil {
			return err
		}
		for _, p := range ps {
			if p.Name == portgroupName {
				portGroupUrn = p.Urn
			}
		}
	}
	vmm := vm.NewManager(c, siteUri)
	res, err := vmm.UploadImage(siteUri+"/vms", vm.ImportTemplateRequest{
		Name:        constant.FusionComputeImageName,
		Description: "KubeOperator默认模版",
		VmConfig: vm.Config{
			Cpu:    vm.Cpu{Quantity: 4, Reservation: 0},
			Memory: vm.Memory{QuantityMB: 4096, Reservation: 4096},
			Disks: []vm.Disk{
				{
					SequenceNum:  1,
					QuantityGB:   50,
					IsDataCopy:   true,
					DatastoreUrn: datastoreUrn,
					IsThin:       true,
				},
			},
			Nics: []vm.Nic{
				{
					Name:         "Network Adapter 0",
					PortGroupUrn: portGroupUrn,
				},
			},
		},
		OsOptions: vm.OsOption{
			OsType:      "Linux",
			OsVersion:   1202,
			GuestOSName: "CentOS 7.6 64bit",
		},
		Location:   clusterUrn,
		Protocol:   "nfs",
		Url:        ovfPath,
		IsTemplate: true,
	})
	if err != nil {
		return err
	}
	tm := task.NewManager(c, siteUri)
	for {
		time.Sleep(5 * time.Second)
		tt, err := tm.Get(res.TaskUri)
		if err != nil {
			return err
		}
		if tt.Status != "running" {
			if tt.Status == "success" {
				break
			} else {
				return errors.New(tt.ReasonDes)
			}
		}
	}
	return nil
}

func (f *fusionComputeClient) DefaultImageExist() (bool, error) {
	siteName, ok := f.Vars["datacenter"].(string)
	if !ok {
		return false, errors.New("type aassertion failed")
	}
	c := f.newFusionComputeClient()
	if err := c.Connect(); err != nil {
		return false, err
	}
	defer func() {
		if err := c.DisConnect(); err != nil {
			log.Errorf("fusionComputeClient DisConnect failed, error: %s", err.Error())
		}
	}()
	sm := site.NewManager(c)
	ss, err := sm.ListSite()
	if err != nil {
		return false, err
	}
	siteUri := ""
	for _, s := range ss {
		if s.Name == siteName {
			siteUri = s.Uri
		}
	}
	if siteUri == "" {
		return false, fmt.Errorf("site %s not found", siteName)
	}
	vmm := vm.NewManager(c, siteUri)
	vms, err := vmm.ListVm(true)
	if err != nil {
		return false, err
	}
	result := false
	for _, tem := range vms {
		if tem.Name == constant.VSphereImageName {
			result = true
			break
		}
	}
	return result, nil
}

func (f *fusionComputeClient) CreateDefaultFolder() error {
	return nil
}

func (f *fusionComputeClient) newFusionComputeClient() client.FusionComputeClient {
	server, ok := f.Vars["server"].(string)
	if !ok {
		log.Errorf("type aassertion failed")
	}
	user, ok := f.Vars["user"].(string)
	if !ok {
		log.Errorf("type aassertion failed")
	}
	password, ok := f.Vars["password"].(string)
	if !ok {
		log.Errorf("type aassertion failed")
	}
	return client.NewFusionComputeClient(server, user, password)
}

func (f *fusionComputeClient) ListDatastores() ([]DatastoreResult, error) {
	var results []DatastoreResult
	siteName, ok := f.Vars["datacenter"].(string)
	if !ok {
		log.Errorf("type aassertion failed")
	}
	c := f.newFusionComputeClient()
	if err := c.Connect(); err != nil {
		return results, err
	}
	defer func() {
		if err := c.DisConnect(); err != nil {
			log.Errorf("fusionComputeClient DisConnect failed, error: %s", err.Error())
		}
	}()
	sm := site.NewManager(c)
	ss, err := sm.ListSite()
	if err != nil {
		return results, err
	}
	siteUri := ""
	for _, s := range ss {
		if s.Name == siteName {
			siteUri = s.Uri
		}
	}
	if siteUri == "" {
		return results, fmt.Errorf("site %s not found", siteName)
	}
	dm := storage.NewManager(c, siteUri)
	datastores, err := dm.ListDataStore()
	if err != nil {
		return results, err
	}
	for i := range datastores {
		results = append(results, DatastoreResult{
			Name:      datastores[i].Name,
			Capacity:  datastores[i].CapacityGB,
			FreeSpace: datastores[i].FreeSizeGB,
		})
	}
	return results, nil
}
