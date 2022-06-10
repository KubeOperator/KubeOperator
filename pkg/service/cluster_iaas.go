package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"

	"github.com/KubeOperator/KubeOperator/pkg/cloud_provider"
	"github.com/KubeOperator/KubeOperator/pkg/cloud_provider/client"
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	"github.com/KubeOperator/KubeOperator/pkg/repository"
	"github.com/KubeOperator/KubeOperator/pkg/util/ipaddr"
	"github.com/KubeOperator/KubeOperator/pkg/util/kotf"
	"github.com/KubeOperator/KubeOperator/pkg/util/lang"
	"github.com/jinzhu/gorm"
)

type ClusterCreateHelper interface {
	LoadMetalNodes(creation *dto.ClusterCreate, cluster *model.Cluster, tx *gorm.DB) error
	LoadPlanNodes(cluster *model.Cluster) error
}

func NewClusterCreateHelper() ClusterCreateHelper {
	return &clusterCreateHelper{
		clusterRepo:         repository.NewClusterRepository(),
		nodeRepo:            repository.NewClusterNodeRepository(),
		hostRepo:            repository.NewHostRepository(),
		planRepo:            repository.NewPlanRepository(),
		projectResourceRepo: repository.NewProjectResourceRepository(),
		vmConfigRepo:        repository.NewVmConfigRepository(),
	}
}

type clusterCreateHelper struct {
	clusterRepo         repository.ClusterRepository
	hostRepo            repository.HostRepository
	nodeRepo            repository.ClusterNodeRepository
	planRepo            repository.PlanRepository
	projectResourceRepo repository.ProjectResourceRepository
	vmConfigRepo        repository.VmConfigRepository
}

func (c clusterCreateHelper) LoadMetalNodes(creation *dto.ClusterCreate, cluster *model.Cluster, tx *gorm.DB) error {
	workerNo, masterNo, firstMasterIP := 1, 1, ""
	for _, nc := range creation.Nodes {
		n := model.ClusterNode{
			ClusterID: cluster.ID,
			Role:      nc.Role,
		}

		var host model.Host
		if err := tx.Set("gorm:query_option", "FOR UPDATE").Where("name = ?", nc.HostName).First(&host).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("can not find host %s cluster id ", nc.HostName)
		}
		host.ClusterID = cluster.ID
		if err := tx.Save(&host).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("can not update host %s cluster id ", nc.HostName)
		}

		clusterResource := model.ClusterResource{
			ClusterID:    cluster.ID,
			ResourceID:   host.ID,
			ResourceType: constant.ResourceHost,
		}
		if err := tx.Create(&clusterResource).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("can bind host %s to cluster", nc.HostName)
		}

		switch cluster.NodeNameRule {
		case constant.NodeNameRuleDefault:
			if n.Role == constant.NodeRoleNameMaster {
				n.Name = fmt.Sprintf("%s-%s-%d", cluster.Name, constant.NodeRoleNameMaster, masterNo)
				if len(firstMasterIP) == 0 {
					firstMasterIP = n.Host.Ip
				}
				masterNo++
			} else {
				n.Name = fmt.Sprintf("%s-%s-%d", cluster.Name, constant.NodeRoleNameWorker, workerNo)
				workerNo++
			}
		case constant.NodeNameRuleIP:
			n.Name = host.Ip
		case constant.NodeNameRuleHostName:
			n.Name = host.Name
		}

		n.HostID = host.ID
		if err := tx.Create(&n).Error; err != nil {
			return fmt.Errorf("can not create  node %s reason %s", n.Name, err.Error())
		}
		n.Host = host
		cluster.Nodes = append(cluster.Nodes, n)
	}
	if cluster.SpecConf.LbMode == constant.LbModeInternal {
		cluster.SpecConf.LbKubeApiserverIp = firstMasterIP
	}
	cluster.SpecConf.KubeRouter = firstMasterIP
	return nil
}

func (c clusterCreateHelper) LoadPlanNodes(cluster *model.Cluster) error {
	if len(cluster.Nodes) > 0 {
		return nil
	}
	tx := db.DB.Begin()
	hosts, err := c.createHosts(*cluster, cluster.Plan)
	if err != nil {
		return err
	}
	if err := tx.Create(&hosts).Error; err != nil {
		tx.Rollback()
		return err
	}

	k := kotf.NewTerraform(&kotf.Config{Cluster: cluster.Name})
	if err = doInit(k, cluster.Plan, hosts); err != nil {
		tx.Rollback()
		return err
	}

	var (
		projectResources []model.ProjectResource
		clusterResources []model.ClusterResource
	)
	prs, err := c.projectResourceRepo.ListByResourceIDAndType(cluster.ID, constant.ResourceCluster)
	if err != nil {
		tx.Rollback()
		return err
	}
	for _, host := range hosts {
		projectResources = append(projectResources, model.ProjectResource{
			ProjectID:    prs[0].ProjectID,
			ResourceID:   host.ID,
			ResourceType: constant.ResourceHost,
		})
		clusterResources = append(clusterResources, model.ClusterResource{
			ClusterID:    cluster.ID,
			ResourceID:   host.ID,
			ResourceType: constant.ResourceHost,
		})
	}
	if err := tx.Create(&projectResources).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Create(&clusterResources).Error; err != nil {
		tx.Rollback()
		return err
	}

	nodes, err := c.createNodes(*cluster, hosts)
	if err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Create(&nodes).Error; err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func (c clusterCreateHelper) createNodes(cluster model.Cluster, hosts []*model.Host) ([]*model.ClusterNode, error) {
	masterNum := 0
	workerNum := 0
	var nodes []*model.ClusterNode
	for _, host := range hosts {
		role := getHostRole(host.Name)
		no := 0
		if role == constant.NodeRoleNameMaster {
			masterNum++
			no = masterNum
		} else {
			workerNum++
			no = workerNum
		}
		node := model.ClusterNode{
			HostID:    host.ID,
			ClusterID: cluster.ID,
			Role:      role,
		}
		if cluster.NodeNameRule == constant.NodeNameRuleDefault {
			node.Name = fmt.Sprintf("%s-%s-%d", cluster.Name, role, no)
		} else {
			node.Name = host.Ip
		}
		nodes = append(nodes, &node)
	}
	if err := c.nodeRepo.BatchSave(nodes); err != nil {
		return nil, err
	}
	return nodes, nil
}

func (c clusterCreateHelper) createHosts(cluster model.Cluster, plan model.Plan) ([]*model.Host, error) {
	var hosts []*model.Host
	masterAmount := 1
	if plan.DeployTemplate != constant.SINGLE {
		masterAmount = 3
	}
	planVars := map[string]string{}
	_ = json.Unmarshal([]byte(plan.Vars), &planVars)

	for i := 0; i < masterAmount; i++ {
		host := model.Host{
			BaseModel: common.BaseModel{},
			Name:      fmt.Sprintf("%s-master-%d", cluster.Name, i+1),
			Port:      22,
			Status:    constant.ClusterCreating,
			ClusterID: cluster.ID,
		}
		if plan.Region.Provider != constant.OpenStack {
			role := getHostRole(host.Name)
			masterConfig, err := c.vmConfigRepo.Get(planVars[fmt.Sprintf("%sModel", role)])
			if err != nil {
				return nil, err
			}
			host.CpuCore = masterConfig.Cpu
			host.Memory = masterConfig.Memory * 1024
		}
		hosts = append(hosts, &host)
	}
	for i := 0; i < cluster.SpecConf.WorkerAmount; i++ {
		host := model.Host{
			BaseModel: common.BaseModel{},
			Name:      fmt.Sprintf("%s-worker-%d", cluster.Name, i+1),
			Port:      22,
			Status:    constant.ClusterCreating,
			ClusterID: cluster.ID,
		}
		if plan.Region.Provider != constant.OpenStack {
			role := getHostRole(host.Name)
			workerConfig, err := c.vmConfigRepo.Get(planVars[fmt.Sprintf("%sModel", role)])
			if err != nil {
				return nil, err
			}
			host.CpuCore = workerConfig.Cpu
			host.Memory = workerConfig.Memory * 1024
		}
		hosts = append(hosts, &host)
	}
	group := allocateZone(plan.Zones, hosts)
	for k, v := range group {
		providerVars := map[string]interface{}{}
		providerVars["provider"] = plan.Region.Provider
		providerVars["datacenter"] = plan.Region.Datacenter
		zoneVars := map[string]interface{}{}
		_ = json.Unmarshal([]byte(k.Vars), &zoneVars)
		providerVars["cluster"] = zoneVars["cluster"]
		_ = json.Unmarshal([]byte(plan.Region.Vars), &providerVars)
		providerVars["datacenter"] = plan.Region.Datacenter
		cloudClient := cloud_provider.NewCloudClient(providerVars)
		err := allocateIpAddr(cloudClient, *k, v, cluster.ID)
		if err != nil {
			return nil, err
		}
		err = allocateDatastore(cloudClient, *k, v)
		if err != nil {
			return nil, err
		}
	}
	return hosts, nil
}

func getHostRole(name string) string {
	if strings.Contains(name, "-master-") {
		return constant.NodeRoleNameMaster
	}
	return constant.NodeRoleNameWorker
}

func doInit(k *kotf.Kotf, plan model.Plan, hosts []*model.Host) error {
	var zonesVars []map[string]interface{}
	for _, zone := range plan.Zones {
		zoneMap := map[string]interface{}{}
		_ = json.Unmarshal([]byte(zone.Vars), &zoneMap)
		zoneMap["key"] = formatZoneName(zone.Name)
		if zoneMap["datastore"] != nil {
			zoneMap["stores"] = formatDatastores(zoneMap["datastore"].([]interface{}))
		}
		zonesVars = append(zonesVars, zoneMap)
	}
	hostsStr, _ := json.Marshal(parseHosts(hosts, plan))
	cloudRegion := map[string]interface{}{
		"datacenter": plan.Region.Datacenter,
		"zones":      zonesVars,
	}
	cloudRegionStr, _ := json.Marshal(&cloudRegion)
	res, err := k.Init(plan.Region.Provider, plan.Region.Vars, string(cloudRegionStr), string(hostsStr))
	if err != nil {
		return err
	}
	if !res.Success {
		return errors.New(res.GetOutput())
	}
	_, err = k.Apply(plan.Region.Vars)
	if err != nil {
		return err
	}
	for i := range hosts {
		hosts[i].Status = constant.ClusterRunning
	}
	return nil
}

func parseHosts(hosts []*model.Host, plan model.Plan) []map[string]interface{} {
	switch plan.Region.Provider {
	case constant.VSphere:
		return parseVsphereHosts(hosts, plan)
	case constant.OpenStack:
		return parseOpenstackHosts(hosts, plan)
	case constant.FusionCompute:
		return parseFusionComputeHosts(hosts, plan)
	}

	return []map[string]interface{}{}
}

func parseVsphereHosts(hosts []*model.Host, plan model.Plan) []map[string]interface{} {
	var results []map[string]interface{}
	for _, h := range hosts {
		var zoneVars map[string]interface{}
		_ = json.Unmarshal([]byte(h.Zone.Vars), &zoneVars)
		zoneVars["key"] = formatZoneName(h.Zone.Name)
		hMap := map[string]interface{}{}
		hMap["name"] = h.Name
		hMap["shortName"] = h.Name
		hMap["cpu"] = h.CpuCore
		hMap["memory"] = h.Memory
		hMap["ip"] = h.Ip
		hMap["zone"] = zoneVars
		hMap["datastore"] = h.Datastore
		hMap["datastoreKey"] = lang.GetStringKey(h.Datastore)
		results = append(results, hMap)
	}
	return results
}

func parseFusionComputeHosts(hosts []*model.Host, plan model.Plan) []map[string]interface{} {
	var results []map[string]interface{}
	for _, h := range hosts {
		var zoneVars map[string]interface{}
		_ = json.Unmarshal([]byte(h.Zone.Vars), &zoneVars)
		zoneVars["key"] = formatZoneName(h.Zone.Name)
		hMap := map[string]interface{}{}
		hMap["name"] = h.Name
		hMap["shortName"] = h.Name
		hMap["cpu"] = h.CpuCore
		hMap["memory"] = h.Memory
		hMap["ip"] = h.Ip
		hMap["zone"] = zoneVars
		hMap["datastore"] = h.Datastore
		results = append(results, hMap)
	}
	return results
}

func parseOpenstackHosts(hosts []*model.Host, plan model.Plan) []map[string]interface{} {
	var results []map[string]interface{}
	planVars := map[string]string{}
	_ = json.Unmarshal([]byte(plan.Vars), &planVars)
	for _, h := range hosts {
		var zoneVars map[string]interface{}
		_ = json.Unmarshal([]byte(h.Zone.Vars), &zoneVars)
		zoneVars["key"] = formatZoneName(h.Zone.Name)
		role := getHostRole(h.Name)
		hMap := map[string]interface{}{}
		hMap["name"] = h.Name
		hMap["shortName"] = h.Name
		hMap["ip"] = h.Ip
		hMap["model"] = planVars[fmt.Sprintf("%sModel", role)]
		hMap["zone"] = zoneVars
		results = append(results, hMap)
	}
	return results
}

func allocateZone(zones []model.Zone, hosts []*model.Host) map[*model.Zone][]*model.Host {
	groupMap := map[*model.Zone][]*model.Host{}
	for i := range hosts {
		hash := i % len(zones)
		groupMap[&zones[hash]] = append(groupMap[&zones[hash]], hosts[i])
		hosts[i].CredentialID = zones[hash].CredentialID
		hosts[i].ZoneID = zones[hash].ID
		hosts[i].Zone = zones[hash]
		var zoneVars map[string]interface{}
		_ = json.Unmarshal([]byte(zones[hash].Vars), &zoneVars)
		if zoneVars["port"] != nil {
			port, _ := strconv.Atoi(zoneVars["port"].(string))
			hosts[i].Port = port
		}
	}
	return groupMap
}

func allocateIpAddr(p cloud_provider.CloudClient, zone model.Zone, hosts []*model.Host, clusterId string) error {
	zoneVars := map[string]string{}
	_ = json.Unmarshal([]byte(zone.Vars), &zoneVars)
	pool, _ := p.GetIpInUsed(zoneVars["network"])
	var hs []model.Host
	if err := db.DB.Find(&hs).Error; err != nil {
		return err
	}
	for i := range hs {
		pool = append(pool, hs[i].Ip)
	}
	var ips []model.Ip
	if err := db.DB.Where("ip_pool_id = ? AND status = ?", zone.IpPoolID, constant.IpAvailable).Order("inet_aton(address)").Find(&ips).Error; err != nil {
		return err
	}
	var wg sync.WaitGroup
	for i := range ips {
		wg.Add(1)
		go func(i int) {
			err := ipaddr.Ping(ips[i].Address)
			if err == nil {
				ips[i].Status = constant.IpReachable
				db.DB.Save(&ips[i])
			}
			wg.Done()
		}(i)
	}
	wg.Wait()

	var aIps []model.Ip
	for i := range ips {
		if ips[i].Status != constant.IpReachable {
			aIps = append(aIps, ips[i])
		}
	}

	var usedIps []model.Ip
	var uIps []string
end:
	for i := range hosts {
		for j := range aIps {
			if !exists(aIps[j].Address, pool) && !exists(aIps[j].Address, uIps) {
				hosts[i].Ip = aIps[j].Address
				usedIps = append(usedIps, aIps[j])
				uIps = append(uIps, aIps[j].Address)
				continue end
			}
		}
	}
	for _, h := range hosts {
		if h.Ip == "" {
			return errors.New("NO_IP_AVAILABLE")
		}
	}
	for i := range usedIps {
		usedIps[i].ClusterID = clusterId
		usedIps[i].Status = constant.IpUsed
		db.DB.Save(&usedIps[i])
	}
	return nil
}

func exists(ip string, pool []string) bool {
	for _, i := range pool {
		if ip == i {
			return true
		}
	}
	return false
}

func formatZoneName(name string) string {
	if lang.CountChinese(name) > 0 {
		return lang.Pinyin(name)
	}
	return name
}

func formatDatastores(datastores []interface{}) map[string]string {
	result := make(map[string]string)
	for _, datastore := range datastores {
		result[datastore.(string)] = lang.GetStringKey(datastore.(string))
	}
	return result
}

func allocateDatastore(p cloud_provider.CloudClient, zone model.Zone, hosts []*model.Host) error {

	zoneVars := map[string]interface{}{}
	_ = json.Unmarshal([]byte(zone.Vars), &zoneVars)
	if zoneVars["datastore"] == nil {
		return nil
	}
	_, ok := zoneVars["datastore"].(string)
	if ok {
		return nil
	}

	var CDatastores []string
	if reflect.TypeOf(zoneVars["datastore"]).Kind() == reflect.Slice {
		s := reflect.ValueOf(zoneVars["datastore"])
		for i := 0; i < s.Len(); i++ {
			ele := s.Index(i)
			CDatastores = append(CDatastores, ele.Interface().(string))
		}
	}

	if len(CDatastores) == 1 {
		for i := range hosts {
			hosts[i].Datastore = CDatastores[0]
		}
		return nil
	}
	results, err := p.ListDatastores()
	if err != nil {
		return err
	}
	var datastores []client.DatastoreResult
	for i := range results {
		for j := range CDatastores {
			if results[i].Name == CDatastores[j] {
				datastores = append(datastores, results[i])
			}
		}
	}

	var chooseDatastore string

	if zoneVars["datastoreType"] == constant.Usage {
		remaining := 0.0
		for i := range datastores {
			dRemaining, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(datastores[i].FreeSpace)/float64(datastores[i].Capacity)), 64)
			if i == 0 {
				remaining = dRemaining
			}
			if dRemaining >= remaining {
				chooseDatastore = datastores[i].Name
			}
		}
	}
	if zoneVars["datastoreType"] == constant.Value {
		value := 0
		for i := range datastores {
			if i == 0 {
				value = datastores[i].FreeSpace
			}
			if datastores[i].FreeSpace >= value {
				chooseDatastore = datastores[i].Name
			}
		}
	}

	for i := range hosts {
		hosts[i].Datastore = chooseDatastore
	}

	return nil
}
