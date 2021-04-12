package client

import (
	"context"
	"errors"
	"io"
	"io/ioutil"
	"net/url"
	"strconv"
	"strings"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/ovf"
	"github.com/vmware/govmomi/property"
	"github.com/vmware/govmomi/view"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/soap"
	"github.com/vmware/govmomi/vim25/types"
)

type vSphereClient struct {
	Vars   map[string]interface{}
	Client *govmomi.Client
}

func NewVSphereClient(vars map[string]interface{}) *vSphereClient {
	return &vSphereClient{
		Vars: vars,
	}
}

func (v *vSphereClient) ListDatacenter() ([]string, error) {
	err := v.GetConnect()
	if err != nil {
		return nil, err
	}
	client := v.Client.Client
	var result []string

	var datacenters []*object.Datacenter
	f := find.NewFinder(client, true)
	datacenters, err = f.DatacenterList(context.TODO(), "*")
	if err != nil {
		return nil, err
	}

	for _, d := range datacenters {
		datacenterPath := d.Common.InventoryPath
		result = append(result, strings.Replace(datacenterPath, "/", "", 1))
	}
	return result, nil
}

func (v *vSphereClient) ListClusters() ([]interface{}, error) {
	err := v.GetConnect()
	if err != nil {
		return nil, err
	}
	client := v.Client.Client
	var result []interface{}

	m := view.NewManager(client)

	view, err := m.CreateContainerView(context.TODO(), client.ServiceContent.RootFolder, []string{"ClusterComputeResource"}, true)
	if err != nil {
		return result, err
	}
	var clusters []mo.ClusterComputeResource
	err = view.Retrieve(context.TODO(), []string{"ClusterComputeResource"}, []string{"summary", "name", "resourcePool", "network", "datastore", "parent"}, &clusters)
	if err != nil {
		return result, err
	}

	pc := property.DefaultCollector(client)
	for _, d := range clusters {

		var host mo.ManagedEntity
		if err := pc.RetrieveOne(context.TODO(), *d.Parent, []string{"name", "parent"}, &host); err != nil {
			return result, err
		}
		var datacenter mo.ManagedEntity
		if err := pc.RetrieveOne(context.TODO(), *host.Parent, []string{"name"}, &datacenter); err != nil {
			return result, err
		}

		if datacenter.Name != v.Vars["datacenter"] {
			continue
		}

		clusterData := make(map[string]interface{})

		clusterData["cluster"] = d.ManagedEntity.Name
		networks, _ := v.GetNetwork(d.ComputeResource.Network)
		clusterData["networks"] = networks
		datastores, _ := v.GetDatastore(d.ComputeResource.Datastore)
		clusterData["datastores"] = datastores
		resourcePools, _ := v.GetResourcePools(*d.ComputeResource.ResourcePool)
		clusterData["resourcePools"] = resourcePools

		result = append(result, clusterData)
	}

	return result, nil
}

func (v *vSphereClient) ListTemplates() ([]interface{}, error) {
	err := v.GetConnect()
	if err != nil {
		return nil, err
	}
	client := v.Client.Client
	var result []interface{}

	m := view.NewManager(client)

	w, err := m.CreateContainerView(context.TODO(), client.ServiceContent.RootFolder, []string{"VirtualMachine"}, true)
	if err != nil {
		return result, err
	}

	var vms []mo.VirtualMachine
	err = w.Retrieve(context.TODO(), []string{"VirtualMachine"}, []string{"summary", "name", "storage"}, &vms)
	if err != nil {
		return result, err
	}

	for _, vm := range vms {
		template := make(map[string]interface{})
		if vm.Summary.Config.Template {
			template["imageName"] = vm.Summary.Config.Name
			template["guestId"] = vm.Summary.Config.GuestId
			var disks []int
			for i := 0; i < int(vm.Summary.Config.NumVirtualDisks); i++ {
				disks = append(disks, i)
			}
			template["imageDisks"] = disks
			result = append(result, template)
		}
	}

	return result, nil
}

func (v *vSphereClient) GetNetwork(mos []types.ManagedObjectReference) ([]string, error) {

	pc := property.DefaultCollector(v.Client.Client)
	rps := []mo.Network{}
	var data []string
	err := pc.Retrieve(context.TODO(), mos, []string{"summary", "name"}, &rps)
	if err != nil {
		return data, err
	}
	for _, d := range rps {
		data = append(data, d.Name)
	}
	return data, nil
}

func (v *vSphereClient) GetDatastore(mos []types.ManagedObjectReference) ([]string, error) {

	pc := property.DefaultCollector(v.Client.Client)
	rps := []mo.Datastore{}
	var data []string
	err := pc.Retrieve(context.TODO(), mos, []string{"summary", "name"}, &rps)
	if err != nil {
		return data, err
	}
	for _, d := range rps {
		data = append(data, d.Name)
	}
	return data, nil
}

func (v *vSphereClient) GetResourcePools(m types.ManagedObjectReference) ([]string, error) {
	pc := property.DefaultCollector(v.Client.Client)
	rp := mo.ResourcePool{}
	var data []string
	err := pc.RetrieveOne(context.TODO(), m, []string{"summary", "name", "resourcePool"}, &rp)
	if err != nil {
		return data, err
	}
	data = append(data, rp.Name)

	rps := []mo.ResourcePool{}
	err = pc.Retrieve(context.TODO(), rp.ResourcePool, []string{"summary", "name"}, &rps)
	if err != nil {
		return data, err
	}
	for _, r := range rps {
		data = append(data, r.Name)
	}

	return data, nil
}

func (v *vSphereClient) GetIpInUsed(network string) ([]string, error) {
	if err := v.GetConnect(); err != nil {
		return nil, err
	}

	var results []string
	c := v.Client.Client
	ctx := context.TODO()
	m := view.NewManager(c)
	vi, err := m.CreateContainerView(ctx, c.ServiceContent.RootFolder, []string{"VirtualMachine", "Network"}, true)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := vi.Destroy(ctx); err != nil {
			log.Errorf("vSphereClient Destroy failed, error: %s", err.Error())
		}
	}()
	var networks []mo.Network
	err = vi.Retrieve(ctx, []string{"Network"}, []string{}, &networks)
	if err != nil {
		return nil, err
	}

	for _, net := range networks {
		if net.Name == network {
			var vms []mo.VirtualMachine
			err = vi.RetrieveWithFilter(ctx, []string{"VirtualMachine"}, []string{"network", "guest"}, &vms, property.Filter{
				"network": net.Reference(),
			})
			if err != nil {
				return nil, err
			}
			for _, vm := range vms {
				for _, n := range vm.Guest.Net {
					results = append(results, n.IpAddress...)
				}
			}
			break
		}
	}
	return results, nil
}

func (v *vSphereClient) GetConnect() error {
	u, err := soap.ParseURL(v.Vars["host"].(string) + ":" + strconv.FormatFloat(v.Vars["port"].(float64), 'G', -1, 64))
	if err != nil {
		return err
	}
	u.User = url.UserPassword(v.Vars["username"].(string), v.Vars["password"].(string))
	c, err := govmomi.NewClient(context.TODO(), u, true)
	if err != nil {
		return err
	}
	v.Client = c
	return nil
}

func (v *vSphereClient) ListFlavors() ([]interface{}, error) {
	return nil, nil
}

func (v *vSphereClient) UploadImage() error {
	if err := v.GetConnect(); err != nil {
		return err
	}
	client := v.Client.Client
	ctx := context.TODO()

	file, _, err := OpenRemoteFile(v.Vars["ovfPath"].(string))
	if err != nil {
		return err
	}
	o, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	f := find.NewFinder(client, true)

	resourcePoolPath := v.Vars["resourcePool"].(string)
	if v.Vars["resourcePool"].(string) == "Resources" {
		resourcePoolPath = "/" + v.Vars["datacenter"].(string) + "/host/" + v.Vars["cluster"].(string) + "/Resources"
	}

	datacenter, err := f.Datacenter(ctx, v.Vars["datacenter"].(string))
	if err != nil {
		return err
	}
	f.SetDatacenter(datacenter)
	resourcePool, err := f.ResourcePool(ctx, resourcePoolPath)
	if err != nil {
		return err
	}
	var datastoreName string
	for _, name := range v.Vars["datastore"].([]interface{}) {
		datastoreName = name.(string)
		break
	}
	datastore, err := f.Datastore(ctx, datastoreName)
	if err != nil {
		return err
	}
	hosts, err := f.HostSystemList(ctx, "*")
	if err != nil {
		return err
	}
	host := hosts[0]

	folder, err := f.Folder(ctx, constant.VSphereFolder)
	if err != nil {
		fd, err := f.DefaultFolder(ctx)
		if err != nil {
			return err
		}
		folder, err = fd.CreateFolder(ctx, constant.VSphereFolder)
		if err != nil {
			return err
		}
	}

	vm, _ := f.VirtualMachine(ctx, constant.VSphereImageName)
	if vm != nil {
		return nil
	}

	var nmap []types.OvfNetworkMapping
	network, err := f.Network(ctx, v.Vars["network"].(string))
	if err != nil {
		return err
	}
	ref := network.Reference()
	nmap = append(nmap, types.OvfNetworkMapping{
		Name:    v.Vars["network"].(string),
		Network: ref,
	})

	cisp := types.OvfCreateImportSpecParams{
		NetworkMapping: nmap,
	}
	ovfClient := ovf.NewManager(client)
	spec, err := ovfClient.CreateImportSpec(ctx, string(o), resourcePool, datastore, cisp)
	if err != nil {
		return err
	}
	if spec.Error != nil {
		return errors.New(spec.Error[0].LocalizedMessage)
	}
	lease, err := resourcePool.ImportVApp(ctx, spec.ImportSpec, folder, host)
	if err != nil {
		return err
	}
	info, err := lease.Wait(ctx, spec.FileItem)
	if err != nil {
		return err
	}
	u := lease.StartUpdater(ctx, info)
	defer u.Done()
	for _, i := range info.Items {
		file, size, err := OpenRemoteFile(v.Vars["vmdkPath"].(string))
		if err != nil {
			return err
		}
		opts := soap.Upload{
			ContentLength: size,
		}
		err = lease.Upload(ctx, i, file, opts)
		if err != nil {
			return err
		}
	}

	err = lease.Complete(ctx)
	if err != nil {
		return err
	}

	template, err := f.VirtualMachine(ctx, constant.VSphereImageName)
	if err != nil {
		return err
	}
	err = template.MarkAsTemplate(ctx)
	if err != nil {
		return err
	}
	return nil
}

func OpenRemoteFile(remoteUrl string) (io.ReadCloser, int64, error) {
	u, err := url.Parse(remoteUrl)
	if err != nil {
		return nil, 0, err
	}
	client := soap.NewClient(u, false)
	f, size, err := client.Download(context.TODO(), u, &soap.DefaultDownload)
	if err != nil {
		return nil, 0, err
	}
	return f, size, nil
}

func (v *vSphereClient) DefaultImageExist() (bool, error) {
	if err := v.GetConnect(); err != nil {
		return false, err
	}
	client := v.Client.Client
	ctx := context.TODO()
	f := find.NewFinder(client, true)
	datacenter, err := f.Datacenter(ctx, v.Vars["datacenter"].(string))
	if err != nil {
		return false, err
	}
	f.SetDatacenter(datacenter)

	vm, err := f.VirtualMachine(ctx, constant.VSphereImageName)
	if err != nil {
		return false, nil
	}
	if vm != nil {
		return true, nil
	}
	return false, nil
}

func (v *vSphereClient) CreateDefaultFolder() error {
	if err := v.GetConnect(); err != nil {
		return err
	}
	client := v.Client.Client
	ctx := context.TODO()
	f := find.NewFinder(client, true)
	datacenter, err := f.Datacenter(ctx, v.Vars["datacenter"].(string))
	if err != nil {
		return err
	}
	f.SetDatacenter(datacenter)

	_, err = f.Folder(ctx, constant.VSphereFolder)
	if err != nil {
		fd, err := f.DefaultFolder(ctx)
		if err != nil {
			return err
		}
		_, err = fd.CreateFolder(ctx, constant.VSphereFolder)
		if err != nil {
			return err
		}
	}

	return nil
}

func (v *vSphereClient) ListDatastores() ([]DatastoreResult, error) {

	var result []DatastoreResult
	if err := v.GetConnect(); err != nil {
		return result, err
	}
	client := v.Client.Client
	ctx := context.TODO()
	m := view.NewManager(client)

	vi, err := m.CreateContainerView(ctx, client.ServiceContent.RootFolder, []string{"ClusterComputeResource"}, true)
	if err != nil {
		return result, err
	}
	defer func() {
		if err := vi.Destroy(ctx); err != nil {
			log.Errorf("vSphereClient Destroy failed, error: %s", err.Error())
		}
	}()
	var clusters []mo.ClusterComputeResource
	err = vi.Retrieve(ctx, []string{"ClusterComputeResource"}, []string{"summary", "name", "resourcePool", "network", "datastore", "parent"}, &clusters)
	if err != nil {
		return result, err
	}
	var dss []mo.Datastore
	for _, d := range clusters {
		if d.Name == v.Vars["cluster"].(string) {
			pc := property.DefaultCollector(v.Client.Client)
			err := pc.Retrieve(ctx, d.ComputeResource.Datastore, []string{"summary", "name"}, &dss)
			if err != nil {
				return result, err
			}
		}
	}

	for i := range dss {
		result = append(result, DatastoreResult{
			Name:      dss[i].Summary.Name,
			Capacity:  int(dss[i].Summary.Capacity / (1024 * 1024 * 1024)),
			FreeSpace: int(dss[i].Summary.FreeSpace / (1024 * 1024 * 1024)),
		})
	}

	return result, nil
}
