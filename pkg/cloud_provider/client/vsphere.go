package client

import (
	"context"
	"errors"
	"github.com/vmware/govmomi/nfc"
	"io"
	"io/ioutil"
	"net/url"
	"strconv"
	"strings"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/logger"
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

func (v *vSphereClient) ListDatacenter() ([]string, string, error) {
	err := v.GetConnect()
	version := ""
	if err != nil {
		return nil, version, err
	}
	client := v.Client.Client
	var result []string
	version = v.Client.ServiceContent.About.Version
	var datacenters []*object.Datacenter
	f := find.NewFinder(client, true)
	datacenters, err = f.DatacenterList(context.TODO(), "*")
	if err != nil {
		return nil, version, err
	}
	for _, d := range datacenters {
		datacenterPath := d.Common.InventoryPath
		result = append(result, strings.Replace(datacenterPath, "/", "", 1))
	}
	return result, version, nil
}

func (v *vSphereClient) ListClusters() ([]interface{}, error) {
	err := v.GetConnect()
	if err != nil {
		return nil, err
	}
	defer v.Client.CloseIdleConnections()
	client := v.Client.Client
	todo := context.TODO()

	f := find.NewFinder(client, true)
	dc, err := f.Datacenter(todo, v.Vars["datacenter"].(string))
	if err != nil {
		return nil, err
	}
	f.SetDatacenter(dc)

	pools, err := f.ResourcePoolListAll(todo, "*")
	if err != nil {
		return nil, err
	}
	resourcePools := make([]string, len(pools))
	for i, value := range pools {
		resourcePools[i] = value.InventoryPath
	}

	dss, err := f.DatastoreList(todo, "*")
	if err != nil {
		return nil, err
	}
	datastores := make([]string, len(dss))
	for i, value := range dss {
		datastores[i] = value.Name()
	}

	nws, err := f.NetworkList(todo, "*")
	if err != nil {
		return nil, err
	}
	networks := []string{}
	for _, value := range nws {
		if value.Reference().Type == "Network" {
			network, _ := v.getNetworkName(value.Reference())
			networks = append(networks, network)
		}
	}

	hs, err := f.HostSystemList(todo, "*")
	if err != nil {
		return nil, err
	}
	hosts := make([]string, len(hs))
	for i, value := range hs {
		hosts[i] = value.Name()
	}

	var result []interface{}
	clusterData := make(map[string]interface{})
	clusterData["networks"] = networks
	clusterData["datastores"] = datastores
	clusterData["resourcePools"] = resourcePools
	clusterData["hosts"] = hosts
	result = append(result, clusterData)

	return result, nil
}

func (v *vSphereClient) getNetworkName(ref types.ManagedObjectReference) (string, error) {
	pc := property.DefaultCollector(v.Client.Client)
	ns := mo.Network{}
	err := pc.RetrieveOne(context.TODO(), ref, []string{"summary", "name"}, &ns)
	if err != nil {
		return "", err
	}
	return ns.Name, nil
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

func (v *vSphereClient) getHosts(mos []types.ManagedObjectReference) ([]interface{}, error) {
	pc := property.DefaultCollector(v.Client.Client)
	var hostResult []interface{}

	var hosts []mo.HostSystem
	if err := pc.Retrieve(context.TODO(), mos, []string{"name", "network", "datastore"}, &hosts); err != nil {
		return nil, err
	}
	for _, h := range hosts {
		hostData := make(map[string]interface{})
		hostData["name"] = h.Name
		hostData["value"] = h.ManagedEntity.Self.Value
		datastores, _ := v.GetDatastore(h.Datastore)
		hostData["datastores"] = datastores
		networks, _ := v.GetNetwork(h.Network)
		hostData["networks"] = networks
		hostResult = append(hostResult, hostData)
	}

	return hostResult, nil
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
			logger.Log.Errorf("vSphereClient Destroy failed, error: %s", err.Error())
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
	defer file.Close()
	o, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	f := find.NewFinder(client, true)

	datacenter, err := f.Datacenter(ctx, v.Vars["datacenter"].(string))
	if err != nil {
		return err
	}
	f.SetDatacenter(datacenter)

	//var resourcePool *object.ResourcePool
	resourcePoolPath := v.Vars["resource"].(string)
	resourcePool, err := f.ResourcePool(ctx, resourcePoolPath)
	if err != nil {
		return err
	}
	var host *object.HostSystem
	hosts, err := f.HostSystemList(ctx, "*")
	if err != nil {
		return err
	}
	resourceType := v.Vars["resourceType"].(string)
	if resourceType == "host" {
		for _, value := range hosts {
			if value.Name() == v.Vars["hostSystem"].(string) {
				host = value
			}
		}
	} else {
		host = hosts[0]
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

	vm, _ := f.VirtualMachine(ctx, v.Vars["imageName"].(string))
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
		EntityName:     v.Vars["imageName"].(string),
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
			_ = CancelUpload(ctx, lease, info.DynamicData)
			return err
		}
		opts := soap.Upload{
			ContentLength: size,
		}
		err = lease.Upload(ctx, i, file, opts)
		if err != nil {
			file.Close()
			_ = CancelUpload(ctx, lease, info.DynamicData)
			return err
		}
		file.Close()
	}

	err = lease.Complete(ctx)
	if err != nil {
		return err
	}

	template, err := f.VirtualMachine(ctx, v.Vars["imageName"].(string))
	if err != nil {
		return err
	}
	err = template.MarkAsTemplate(ctx)
	if err != nil {
		return err
	}
	return nil
}

func CancelUpload(ctx context.Context, lease *nfc.Lease, data types.DynamicData) error {
	return lease.Abort(ctx, &types.LocalizedMethodFault{
		DynamicData: data,
		Fault: &types.OvfImportFailed{
			types.OvfImport{},
		},
	})
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

func (v *vSphereClient) ImageExist(template string) (bool, error) {
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

	vm, err := f.VirtualMachine(ctx, template)
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

	f := find.NewFinder(client, true)
	dc, err := f.Datacenter(ctx, v.Vars["datacenter"].(string))
	if err != nil {
		return nil, err
	}
	f.SetDatacenter(dc)

	dss, err := f.DatastoreList(ctx, "*")
	if err != nil {
		return nil, err
	}
	var objs []types.ManagedObjectReference

	for _, ds := range dss {
		objs = append(objs, ds.Reference())
	}
	var datastores []mo.Datastore
	pc := property.DefaultCollector(v.Client.Client)
	err = pc.Retrieve(ctx, objs, []string{"summary", "name"}, &datastores)
	if err != nil {
		return result, err
	}

	for i := range datastores {
		result = append(result, DatastoreResult{
			Name:      datastores[i].Summary.Name,
			Capacity:  int(datastores[i].Summary.Capacity / (1024 * 1024 * 1024)),
			FreeSpace: int(datastores[i].Summary.FreeSpace / (1024 * 1024 * 1024)),
		})
	}

	return result, nil
}

func (v *vSphereClient) ListFolders() ([]string, error) {
	result := []string{}
	if err := v.GetConnect(); err != nil {
		return result, err
	}
	client := v.Client.Client
	ctx := context.TODO()
	f := find.NewFinder(client, true)
	dc, err := f.Datacenter(ctx, v.Vars["datacenter"].(string))
	if err != nil {
		return nil, err
	}
	f.SetDatacenter(dc)
	folders, err := f.FolderList(ctx, "*")
	if err != nil {
		return nil, err
	}
	for _, fold := range folders {
		result = append(result, fold.Name())
	}
	return result, nil
}
