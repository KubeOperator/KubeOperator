package client

import (
	"encoding/json"
	"errors"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"io/ioutil"
	"net/http"
)

var (
	GetRegionError = "GET_REGION_ERROR"
)

type openStackClient struct {
	Vars map[string]interface{}
}

func NewOpenStackClient(vars map[string]interface{}) *openStackClient {
	return &openStackClient{
		Vars: vars,
	}
}

func (v *openStackClient) ListZones() string {
	return ""
}

func (v *openStackClient) ListDatacenter() ([]string, error) {
	var result []string

	provider, err := v.GetAuth()
	if err != nil {
		return result, err
	}

	client := &http.Client{}
	req, _ := http.NewRequest("GET", v.Vars["identity"].(string)+"/regions", nil)
	req.Header.Add("X-Auth-Token", provider.TokenID)
	resp, _ := client.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)
	m := make(map[string]interface{})
	json.Unmarshal([]byte(body), &m)
	key, exist := m["regions"]
	if exist {
		regions := key.([]interface{})
		for _, r := range regions {
			region := r.(map[string]interface{})
			result = append(result, region["id"].(string))
		}
	} else {
		return result, errors.New(GetRegionError)
	}

	return result, nil
}

func (v *openStackClient) ListClusters() ([]interface{}, error) {
	return []interface{}{}, nil
}
func (v *openStackClient) ListTemplates() ([]interface{}, error) {
	return []interface{}{}, nil
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
