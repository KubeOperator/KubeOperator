package client

import (
	"context"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/soap"
	"net/url"
	"strings"
)

type vSphereClient struct {
	Vars    map[string]interface{}
	Connect Connect
}

type Connect struct {
	Client govmomi.Client
	Ctx    context.Context
}

func NewVSphereClient(vars map[string]interface{}) *vSphereClient {
	return &vSphereClient{
		Vars: vars,
	}
}

func (v *vSphereClient) ListZones() string {
	return ""
}

func (v *vSphereClient) ListDatacenter() ([]string, error) {
	_, err := v.GetConnect()
	if err != nil {
		return nil, err
	}
	client := v.Connect.Client.Client
	//m := view.NewManager(client)
	var data []string
	//view, err := m.CreateContainerView(v.Connect.Ctx, client.ServiceContent.RootFolder, []string{"Datastore"}, true)
	//if err != nil {
	//	return data, err
	//}
	//var datacenters []mo.Datastore
	//err = view.Retrieve(v.Connect.Ctx, []string{"Datastore"}, []string{"summary"}, &datacenters)
	//if err != nil {
	//	return data, err
	//}
	var datacenters []*object.Datacenter
	f := find.NewFinder(client, true)
	datacenters, err = f.DatacenterList(v.Connect.Ctx, "*")
	if err != nil {
		return nil, err
	}

	for _, d := range datacenters {
		data = append(data, d.Common.InventoryPath)
	}
	return data, nil
}

func (v *vSphereClient) ListClusters(datacenter string) ([]string, error) {
	_, err := v.GetConnect()
	if err != nil {
		return nil, err
	}
	client := v.Connect.Client.Client
	//m := view.NewManager(client)
	var data []string
	//view, err := m.CreateContainerView(v.Connect.Ctx, client.ServiceContent.RootFolder, []string{"ClusterComputeResource"}, true)
	//if err != nil {
	//	return data,err
	//}
	//
	//var clusters []object.ClusterComputeResource
	//err = view.Retrieve(v.Connect.Ctx, []string{"ClusterComputeResource"}, []string{"summary"}, &clusters)
	//if err != nil {
	//	return data,err
	//}

	var clusters []*object.ComputeResource
	var dc *object.Datacenter
	f := find.NewFinder(client, true)
	dc, _ = f.Datacenter(v.Connect.Ctx, datacenter)
	f.SetDatacenter(dc)
	clusters, err = f.ComputeResourceList(v.Connect.Ctx, "*")
	if err != nil {
		return nil, err
	}

	for _, d := range clusters {
		var name string
		name = strings.Replace(d.Common.InventoryPath, datacenter+"/host/", "", -1)
		data = append(data, name)
	}
	return data, nil
}

func (v *vSphereClient) GetConnect() (Connect, error) {
	ctx, _ := context.WithCancel(context.Background())
	u, err := soap.ParseURL(v.Vars["vcHost"].(string))
	if err != nil {
		return Connect{}, err
	}
	u.User = url.UserPassword(v.Vars["vcUsername"].(string), v.Vars["vcPassword"].(string))
	c, err := govmomi.NewClient(ctx, u, true)
	if err != nil {
		return Connect{}, err
	}
	connect := &Connect{
		Client: *c,
		Ctx:    ctx,
	}
	v.Connect = *connect
	return *connect, nil
}
