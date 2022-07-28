package client

import (
	"encoding/json"
	"errors"
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/v3/volumetypes"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/availabilityzones"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/secgroups"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/flavors"
	"github.com/gophercloud/gophercloud/openstack/imageservice/v2/imagedata"
	"github.com/gophercloud/gophercloud/openstack/imageservice/v2/imageimport"
	"github.com/gophercloud/gophercloud/openstack/imageservice/v2/images"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/external"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/vlantransparent"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/networks"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/subnets"
	uuid "github.com/satori/go.uuid"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

var (
	GetRegionError = "GET_REGION_ERROR"
)

type networkWithExternalExt struct {
	networks.Network
	external.NetworkExternalExt
	vlantransparent.TransparentExt
}

type openStackClient struct {
	Vars map[string]interface{}
}

func NewOpenStackClient(vars map[string]interface{}) *openStackClient {
	return &openStackClient{
		Vars: vars,
	}
}

func (v *openStackClient) ListDatacenter() ([]string, string, error) {
	var result []string
	version := ""

	provider, err := v.GetAuth()
	if err != nil {
		return result, version, err
	}

	client := &http.Client{}
	req, _ := http.NewRequest("GET", v.Vars["identity"].(string)+"/regions", nil)
	req.Header.Add("X-Auth-Token", provider.TokenID)
	resp, _ := client.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(body), &m); err != nil {
		return result, version, err
	}
	key, exist := m["regions"]
	if exist {
		regions := key.([]interface{})
		for _, r := range regions {
			region := r.(map[string]interface{})
			result = append(result, region["id"].(string))
		}
	} else {
		return result, version, errors.New(GetRegionError)
	}

	return result, version, nil
}

func (v *openStackClient) ListClusters() ([]interface{}, error) {
	var result []interface{}

	provider, err := v.GetAuth()
	if err != nil {
		return result, err
	}

	client, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{
		Region: v.Vars["datacenter"].(string),
	})
	if err != nil {
		return result, err
	}

	pager, err := availabilityzones.List(client).AllPages()
	if err != nil {
		return result, err
	}
	zones, err := availabilityzones.ExtractAvailabilityZones(pager)
	if err != nil {
		return result, err
	}

	sgPages, err := secgroups.List(client).AllPages()
	if err != nil {
		return result, err
	}
	allSecurityGroups, err := secgroups.ExtractSecurityGroups(sgPages)
	if err != nil {
		return result, err
	}

	iPages, err := images.List(client, images.ListOpts{}).AllPages()
	if err != nil {
		return result, err
	}

	allImages, err := images.ExtractImages(iPages)
	if err != nil {
		panic(err)
	}

	networkClient, err := openstack.NewNetworkV2(provider, gophercloud.EndpointOpts{
		Region: v.Vars["datacenter"].(string),
	})
	if err != nil {
		return result, err
	}

	networkPager, err := networks.List(networkClient, networks.ListOpts{}).AllPages()
	if err != nil {
		return result, err
	}

	var allNetworks []networkWithExternalExt
	err = networks.ExtractNetworksInto(networkPager, &allNetworks)
	if err != nil {
		return result, err
	}
	blockStorageClient, err := openstack.NewBlockStorageV3(provider, gophercloud.EndpointOpts{
		Region: v.Vars["datacenter"].(string),
	})
	if err != nil {
		return result, err
	}

	vPages, err := volumetypes.List(blockStorageClient, volumetypes.ListOpts{}).AllPages()
	if err != nil {
		return result, err
	}
	allVPages, err := volumetypes.ExtractVolumeTypes(vPages)
	if err != nil {
		return result, err
	}

	var ipTypes []string
	ipTypes = append(ipTypes, "private")
	ipTypes = append(ipTypes, "floating")

	for _, z := range zones {
		clusterData := make(map[string]interface{})
		clusterData["cluster"] = z.ZoneName

		var networkList []interface{}
		var floatingNetworkList []interface{}

		for _, n := range allNetworks {
			networkData := make(map[string]interface{})
			networkData["name"] = n.Name
			networkData["id"] = n.ID
			subnetPages, err := subnets.List(networkClient, subnets.ListOpts{
				NetworkID: n.ID,
			}).AllPages()
			if err != nil {
				continue
			}
			allSubnets, err := subnets.ExtractSubnets(subnetPages)
			if err != nil {
				continue
			}
			var subnetList []interface{}
			for _, s := range allSubnets {
				subnetData := make(map[string]interface{})
				subnetData["id"] = s.ID
				subnetData["name"] = s.Name
				subnetList = append(subnetList, subnetData)
			}
			networkData["subnetList"] = subnetList
			networkList = append(networkList, networkData)

			if n.NetworkExternalExt.External {
				floatingNetworkList = append(floatingNetworkList, networkData)
			}
		}
		clusterData["networkList"] = networkList
		clusterData["floatingNetworkList"] = floatingNetworkList

		var securityGroups []string
		for _, s := range allSecurityGroups {
			securityGroups = append(securityGroups, s.Name)
		}
		clusterData["securityGroups"] = securityGroups

		var volumeTypes []interface{}
		for _, d := range allVPages {
			volumeData := make(map[string]interface{})
			volumeData["name"] = d.Name
			volumeData["id"] = d.ID
			volumeTypes = append(volumeTypes, volumeData)
		}
		clusterData["storages"] = volumeTypes

		var imageList []interface{}
		for _, i := range allImages {
			imageData := make(map[string]interface{})
			imageData["name"] = i.Name
			imageData["id"] = i.ID
			imageList = append(imageList, imageData)
		}
		clusterData["imageList"] = imageList
		clusterData["ipTypes"] = ipTypes

		result = append(result, clusterData)
	}

	return result, nil
}
func (v *openStackClient) ListTemplates() ([]interface{}, error) {

	var result []interface{}
	provider, err := v.GetAuth()
	if err != nil {
		return result, err
	}

	client, err := openstack.NewImageServiceV2(provider, gophercloud.EndpointOpts{
		Region: v.Vars["datacenter"].(string),
	})
	if err != nil {
		return result, err
	}

	pager, err := images.List(client, images.ListOpts{}).AllPages()
	if err != nil {
		return result, err
	}
	allPages, err := images.ExtractImages(pager)
	if err != nil {
		return result, err
	}

	for _, p := range allPages {
		template := make(map[string]string)
		template["imageName"] = p.Name
		template["id"] = p.ID
		result = append(result, template)
	}
	return result, nil
}

func (v *openStackClient) GetAuth() (*gophercloud.ProviderClient, error) {

	scope := gophercloud.AuthScope{
		ProjectID: v.Vars["projectId"].(string),
	}

	opts := gophercloud.AuthOptions{
		IdentityEndpoint: v.Vars["identity"].(string),
		Username:         v.Vars["username"].(string),
		Password:         v.Vars["password"].(string),
		DomainName:       v.Vars["domainName"].(string),
		Scope:            &scope,
	}

	provider, err := openstack.AuthenticatedClient(opts)
	if err != nil {
		return nil, err
	}
	return provider, nil
}

func (v *openStackClient) ListFlavors() ([]interface{}, error) {

	var result []interface{}

	provider, err := v.GetAuth()
	if err != nil {
		return result, err
	}

	client, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{
		Region: v.Vars["datacenter"].(string),
	})
	if err != nil {
		return result, err
	}
	pager, err := flavors.ListDetail(client, flavors.ListOpts{}).AllPages()
	if err != nil {
		return result, err
	}
	allPages, err := flavors.ExtractFlavors(pager)
	if err != nil {
		return result, err
	}

	for _, f := range allPages {

		if f.RAM > 1024 {
			vmConfig := make(map[string]interface{})
			vmConfig["name"] = f.Name

			config := make(map[string]interface{})
			config["id"], _ = strconv.Atoi(f.ID)
			config["disk"] = f.Disk
			config["cpu"] = f.VCPUs
			config["memory"] = f.RAM / 1024

			vmConfig["config"] = config
			result = append(result, vmConfig)
		}
	}

	return result, nil
}

func (v *openStackClient) GetIpInUsed(network string) ([]string, error) {
	return []string{}, nil
}
func (v *openStackClient) UploadImage() error {

	provider, err := v.GetAuth()
	if err != nil {
		return err
	}

	client, err := openstack.NewImageServiceV2(provider, gophercloud.EndpointOpts{
		Region: v.Vars["datacenter"].(string),
	})
	if err != nil {
		return err
	}

	pager, err := images.List(client, images.ListOpts{}).AllPages()
	if err != nil {
		return err
	}
	allPages, err := images.ExtractImages(pager)
	if err != nil {
		return err
	}

	exist := false
	for _, p := range allPages {
		if p.Name == constant.OpenStackImageName {
			exist = true
			break
		}
	}
	if !exist {
		imageId := uuid.NewV4().String()
		//download image
		res, err := http.Get(v.Vars["imagePath"].(string))
		if err != nil {
			return err
		}
		f, err := os.Create(constant.OpenStackImageLocalPath)
		if err != nil {
			return err
		}
		_, err = io.Copy(f, res.Body)
		if err != nil {
			return err
		}

		create := images.Create(client, images.CreateOpts{
			Name:            constant.OpenStackImageName,
			DiskFormat:      constant.OpenStackImageDiskFormat,
			ContainerFormat: "bare",
			ID:              imageId,
		})
		if create.Err != nil {
			return create.Err
		}

		imageData, err := os.Open(constant.OpenStackImageLocalPath)
		if err != nil {
			return err
		}
		defer imageData.Close()
		err = imagedata.Stage(client, imageId, imageData).ExtractErr()
		if err != nil {
			return err
		}

		result := imageimport.Create(client, imageId, imageimport.CreateOpts{
			Name: imageimport.GlanceDirectMethod,
		})
		if result.Err != nil {
			return result.Err
		}
	}

	return nil
}

func (v *openStackClient) ImageExist(template string) (bool, error) {
	return false, nil
}

func (v *openStackClient) CreateDefaultFolder() error {
	return nil
}

func (v *openStackClient) ListDatastores() ([]DatastoreResult, error) {
	var results []DatastoreResult
	return results, nil
}

func (v *openStackClient) ListFolders() ([]string, error) {
	folders := []string{}
	return folders, nil
}
