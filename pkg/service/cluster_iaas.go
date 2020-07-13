package service

import (
	"fmt"
	"github.com/KubeOperator/KubeOperator/pkg/cloud_provider/client"
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	"github.com/KubeOperator/KubeOperator/pkg/repository"
	"github.com/KubeOperator/KubeOperator/pkg/util/ipaddr"
)

type ClusterIaasService interface {
	Init(name string) error
}

type clusterIaasService struct {
	ClusterService ClusterService
	regionService  repository.RegionRepository
}

//func (c clusterIaasService) Init(name string) error {
//	//cluster, err := c.ClusterService.Get(name)
//	//if err != nil {
//	//	return err
//	//}
//	//plan, err := c.ClusterService.GetPlan(name)
//	//if err != nil {
//	//	return err
//	//}
//}

func (c clusterIaasService) createHosts(cluster model.Cluster, plan model.Plan) error {
	var hosts []*model.Host
	masterAmount := 1
	if plan.DeployTemplate != constant.SINGLE {
		masterAmount = 3
	}
	for i := 1; i < masterAmount; i++ {
		host := model.Host{
			BaseModel: common.BaseModel{},
			Name:      fmt.Sprintf("%s-master-%d", cluster.Name, i),
			Port:      22,
			Status:    constant.ClusterWaiting,
		}
		hosts = append(hosts, &host)
	}
	for i := 1; i < cluster.Spec.WorkerAmount; i++ {
		host := model.Host{
			BaseModel: common.BaseModel{},
			Name:      fmt.Sprintf("%s-worker-%d", cluster.Name, i),
			Port:      22,
			Status:    constant.ClusterWaiting,
		}
		hosts = append(hosts, &host)
	}
	var selectedIps []string
	group := allocateZone(plan.Zones, hosts)
	for k, v := range group {
		cloudClient := client.NewCloudClient(map[string]interface{}{})
		err := allocateIpAddr(cloudClient, *k, v, selectedIps)
		if err != nil {
			return err
		}
	}
	return nil
}

func allocateZone(zones []model.Zone, hosts []*model.Host) map[*model.Zone][]*model.Host {
	groupMap := map[*model.Zone][]*model.Host{}
	for i, _ := range hosts {
		hash := i % len(zones)
		groupMap[&zones[hash]] = append(groupMap[&zones[hash]], hosts[i])
	}
	return groupMap
}

func allocateIpAddr(p client.CloudClient, zone model.Zone, hosts []*model.Host, selectedIps []string) error {
	ips := ipaddr.GenerateIps("172.16.10.0", 24)
	pool, err := p.GetIpInUsed("")
	if err != nil {
		return err
	}
	for _, h := range hosts {
		for i, _ := range ips {
			if !exists(ips[i], pool) && !exists(ips[i], selectedIps) {
				h.Ip = ips[i]
			}
		}
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
