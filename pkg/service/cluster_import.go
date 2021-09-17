package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/logger"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/repository"
	kubeUtil "github.com/KubeOperator/KubeOperator/pkg/util/kubernetes"
	"github.com/icza/dyno"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type ClusterImportService interface {
	Import(clusterImport dto.ClusterImport) error
	LoadClusterInfo(loadInfo *dto.ClusterLoad) (dto.ClusterLoadInfo, error)
}

type clusterImportService struct {
	clusterRepo               repository.ClusterRepository
	projectRepository         repository.ProjectRepository
	projectResourceRepository repository.ProjectResourceRepository
}

func NewClusterImportService() *clusterImportService {
	return &clusterImportService{
		clusterRepo:               repository.NewClusterRepository(),
		projectRepository:         repository.NewProjectRepository(),
		projectResourceRepository: repository.NewProjectResourceRepository(),
	}
}

func (c clusterImportService) Import(clusterImport dto.ClusterImport) error {
	loginfo, _ := json.Marshal(clusterImport)
	logger.Log.WithFields(logrus.Fields{"cluster_import_info": string(loginfo)}).Debugf("start to import the cluster %s", clusterImport.Name)

	project, err := c.projectRepository.Get(clusterImport.ProjectName)
	if err != nil {
		return err
	}

	var address string
	var port int
	if strings.HasSuffix(clusterImport.ApiServer, "/") {
		clusterImport.ApiServer = strings.Replace(clusterImport.ApiServer, "/", "", -1)
	}
	clusterImport.ApiServer = strings.Replace(clusterImport.ApiServer, "http://", "", -1)
	clusterImport.ApiServer = strings.Replace(clusterImport.ApiServer, "https://", "", -1)
	if strings.Contains(clusterImport.ApiServer, ":") {
		strs := strings.Split(clusterImport.ApiServer, ":")
		address = strs[0]
		port, _ = strconv.Atoi(strs[1])
	} else {
		address = clusterImport.ApiServer
		port = 80
	}
	tx := db.DB.Begin()
	cluster := model.Cluster{
		Name:      clusterImport.Name,
		ProjectID: project.ID,
		Source:    constant.ClusterSourceExternal,
		Status: model.ClusterStatus{
			Phase: constant.ClusterRunning,
		},
		Spec: model.ClusterSpec{
			LbKubeApiserverIp: address,
			KubeApiServerPort: port,
			Architectures:     clusterImport.Architectures,
			KubeRouter:        clusterImport.Router,
		},
		Secret: model.ClusterSecret{
			KubeadmToken:    "",
			KubernetesToken: clusterImport.Token,
		},
	}
	if clusterImport.IsKoCluster {
		cluster.Source = constant.ClusterSourceKoExternal
		cluster.Spec = model.ClusterSpec{
			RuntimeType:              clusterImport.KoClusterInfo.RuntimeType,
			DockerStorageDir:         clusterImport.KoClusterInfo.DockerStorageDIr,
			ContainerdStorageDir:     clusterImport.KoClusterInfo.ContainerdStorageDIr,
			NetworkType:              clusterImport.KoClusterInfo.NetworkType,
			CiliumVersion:            clusterImport.KoClusterInfo.CiliumVersion,
			CiliumTunnelMode:         clusterImport.KoClusterInfo.CiliumTunnelMode,
			CiliumNativeRoutingCidr:  clusterImport.KoClusterInfo.CiliumNativeRoutingCidr,
			Version:                  clusterImport.KoClusterInfo.Version,
			Provider:                 constant.ClusterProviderBareMetal,
			FlannelBackend:           clusterImport.KoClusterInfo.FlannelBackend,
			CalicoIpv4poolIpip:       clusterImport.KoClusterInfo.CalicoIpv4poolIpip,
			KubeProxyMode:            clusterImport.KoClusterInfo.KubeProxyMode,
			NodeportAddress:          clusterImport.KoClusterInfo.NodeportAddress,
			KubeServiceNodePortRange: clusterImport.KoClusterInfo.KubeServiceNodePortRange,
			EnableDnsCache:           clusterImport.KoClusterInfo.EnableDnsCache,
			DnsCacheVersion:          clusterImport.KoClusterInfo.DnsCacheVersion,
			IngressControllerType:    clusterImport.KoClusterInfo.IngressControllerType,
			KubernetesAudit:          clusterImport.KoClusterInfo.KubernetesAudit,
			DockerSubnet:             clusterImport.KoClusterInfo.DockerSubnet,
			HelmVersion:              clusterImport.KoClusterInfo.HelmVersion,
			NetworkInterface:         clusterImport.KoClusterInfo.NetworkInterface,
			NetworkCidr:              clusterImport.KoClusterInfo.NetworkCidr,
			SupportGpu:               clusterImport.KoClusterInfo.SupportGpu,
			YumOperate:               clusterImport.KoClusterInfo.YumOperate,

			LbMode:            constant.ClusterSourceInternal,
			LbKubeApiserverIp: address,
			KubeApiServerPort: port,
			Architectures:     clusterImport.Architectures,
			KubeRouter:        clusterImport.Router,

			KubePodSubnet:         clusterImport.KoClusterInfo.KubePodSubnet,
			KubeServiceSubnet:     clusterImport.KoClusterInfo.KubeServiceSubnet,
			MaxNodeNum:            clusterImport.KoClusterInfo.MaxNodeNum,
			KubeMaxPods:           clusterImport.KoClusterInfo.KubeMaxPods,
			KubeNetworkNodePrefix: clusterImport.KoClusterInfo.KubeNetworkNodePrefix,
		}
	}
	if err := tx.Create(&cluster.Spec).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("can not create cluster spec %s", err.Error())
	}
	if err := tx.Create(&cluster.Status).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("can not create cluster status %s", err.Error())
	}
	if err := tx.Create(&cluster.Secret).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("can not create cluster secret %s", err.Error())
	}
	cluster.SpecID = cluster.Spec.ID
	cluster.StatusID = cluster.Status.ID
	cluster.SecretID = cluster.Secret.ID
	if err := tx.Create(&cluster).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("can not create cluster secret %s", err.Error())
	}

	if err := tx.Save(&cluster.Spec).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("can not update spec %s", err.Error())
	}

	var synchosts []dto.HostSync
	if clusterImport.IsKoCluster {
		for _, node := range clusterImport.KoClusterInfo.Nodes {
			host := model.Host{
				Name:         node.Name,
				Ip:           node.Ip,
				Port:         node.Port,
				CredentialID: node.CredentialID,
				ClusterID:    cluster.ID,
				Status:       constant.StatusInitializing,
				Architecture: node.Architecture,
			}
			if err := tx.Create(&host).Error; err != nil {
				tx.Rollback()
				return err
			}
			synchosts = append(synchosts, dto.HostSync{HostName: node.Name, HostStatus: constant.StatusRunning})
			node := model.ClusterNode{
				Name:      node.Name,
				HostID:    host.ID,
				ClusterID: cluster.ID,
				Role:      node.Role,
				Status:    constant.StatusRunning,
			}
			if err := tx.Create(&node).Error; err != nil {
				tx.Rollback()
				return err
			}
			clusterResource := model.ClusterResource{
				ResourceType: constant.ResourceHost,
				ResourceID:   host.ID,
				ClusterID:    cluster.ID,
			}
			if err := tx.Create(&clusterResource).Error; err != nil {
				tx.Rollback()
				return err
			}
			projectResource := model.ProjectResource{
				ResourceType: constant.ResourceHost,
				ResourceID:   host.ID,
				ProjectID:    project.ID,
			}
			if err := tx.Create(&projectResource).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	} else {
		if err := gatherClusterInfo(&cluster); err != nil {
			tx.Rollback()
			return err
		}
		for _, node := range cluster.Nodes {
			node.ClusterID = cluster.ID
			if err := tx.Create(&node).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("can not save node %s", err.Error())
			}
		}
	}

	var (
		manifest model.ClusterManifest
		toolVars []model.VersionHelp
	)
	if err := tx.Where("name = ?", cluster.Spec.Version).Order("created_at ASC").First(&manifest).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("can find manifest version: %s", err.Error())
	}
	if err := json.Unmarshal([]byte(manifest.ToolVars), &toolVars); err != nil {
		tx.Rollback()
		return fmt.Errorf("unmarshal manifest.toolvar error %s", err.Error())
	}
	for _, tool := range cluster.PrepareTools() {
		for _, item := range toolVars {
			if tool.Name == item.Name {
				tool.Version = item.Version
				break
			}
		}
		tool.ClusterID = cluster.ID
		if err := tx.Create(&tool).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("can not save tool %s", err.Error())
		}
	}
	istios := cluster.PrepareIstios()
	for _, istio := range istios {
		istio.ClusterID = cluster.ID
		if err := tx.Create(&istio).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("can not save istio %s", err.Error())
		}
	}
	if err := c.projectResourceRepository.Create(model.ProjectResource{
		ResourceID:   cluster.ID,
		ProjectID:    project.ID,
		ResourceType: constant.ResourceCluster,
	}); err != nil {
		tx.Rollback()
		return fmt.Errorf("can not create project resource %s", err.Error())
	}
	tx.Commit()

	hostService := NewHostService()
	go hostService.SyncList(synchosts)
	return nil
}

func gatherClusterInfo(cluster *model.Cluster) error {
	c, err := kubeUtil.NewKubernetesClient(&kubeUtil.Config{
		Hosts: []kubeUtil.Host{kubeUtil.Host(fmt.Sprintf("%s:%d", cluster.Spec.LbKubeApiserverIp, cluster.Spec.KubeApiServerPort))},
		Token: cluster.Secret.KubernetesToken,
	})
	if err != nil {
		return err
	}
	_, err = c.ServerVersion()
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	for _, f := range funcList {
		wg.Add(1)
		go f(cluster, c, &wg)
	}
	wg.Wait()
	return nil
}

type GatherClusterInfoFunc func(cluster *model.Cluster, client *kubernetes.Clientset, wg *sync.WaitGroup)

var funcList = []GatherClusterInfoFunc{
	getServerVersion,
	getKubeNodes,
	getNetworkType,
	getRuntimeType,
}

func getServerVersion(cluster *model.Cluster, client *kubernetes.Clientset, wg *sync.WaitGroup) {
	defer wg.Done()
	v, err := client.ServerVersion()
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}
	var manifest model.ClusterManifest
	if err := db.DB.Where("version = ?", v.GitVersion).Order("created_at ASC").First(&manifest).Error; err != nil {
		logger.Log.Error("get manifest %s failed: %v", v.GitVersion, err.Error())
	}
	cluster.Spec.Version = manifest.Name
}

func getKubeNodes(cluster *model.Cluster, client *kubernetes.Clientset, wg *sync.WaitGroup) {
	defer wg.Done()
	nodes, err := client.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}
	for _, node := range nodes.Items {
		var role string
		_, ok := node.Labels["node-role.kubernetes.io/master"]
		if ok {
			role = constant.NodeRoleNameMaster
		} else {
			_, ok := node.Labels["node-role.kubernetes.io/worker"]
			if ok {
				role = constant.NodeRoleNameWorker
			}
		}
		cluster.Nodes = append(cluster.Nodes, model.ClusterNode{
			Name:   node.Name,
			Role:   role,
			Status: constant.ClusterRunning,
		})
	}
}

func getNetworkType(cluster *model.Cluster, client *kubernetes.Clientset, wg *sync.WaitGroup) {
	defer wg.Done()
	dps, err := client.AppsV1().DaemonSets("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}
	networkMap := map[string]int{
		"flannel": 0,
		"calico":  0,
	}
	for _, dp := range dps.Items {
		for i := range networkMap {
			if strings.Contains(dp.Name, i) {
				networkMap[i]++
			}
		}
	}
	var networkType = ""
	for k, v := range networkMap {
		if v > 0 {
			networkType = k
			break
		}
	}
	if networkType == "" {
		networkType = "unknown"
	}
	cluster.Spec.NetworkType = networkType
}

func getRuntimeType(cluster *model.Cluster, client *kubernetes.Clientset, wg *sync.WaitGroup) {
	defer wg.Done()
	ns, err := client.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}
	var node v1.Node
	// 找第一个 master
	for _, n := range ns.Items {
		if _, ok := n.Labels["node-role.kubernetes.io/master"]; ok {
			node = n
			break
		}
	}
	cluster.Spec.Architectures = node.Status.NodeInfo.Architecture
	if strings.Contains(node.Status.NodeInfo.ContainerRuntimeVersion, "docker") {
		cluster.Spec.RuntimeType = "docker"
	}
	if strings.Contains(node.Status.NodeInfo.ContainerRuntimeVersion, "containerd") {
		cluster.Spec.RuntimeType = "containerd"
	}
}

func (c clusterImportService) LoadClusterInfo(loadInfo *dto.ClusterLoad) (dto.ClusterLoadInfo, error) {
	var clusterInfo dto.ClusterLoadInfo
	loadInfo.ApiServer = strings.Replace(loadInfo.ApiServer, "http://", "", -1)
	loadInfo.ApiServer = strings.Replace(loadInfo.ApiServer, "https://", "", -1)
	clusterInfo.Architectures = loadInfo.Architectures

	kubeClient, err := kubeUtil.NewKubernetesClient(&kubeUtil.Config{
		Hosts: []kubeUtil.Host{kubeUtil.Host(loadInfo.ApiServer)},
		Token: loadInfo.Token,
	})
	if err != nil {
		return clusterInfo, err
	}

	// load version
	v, err := kubeClient.ServerVersion()
	if err != nil {
		return clusterInfo, err
	}
	var manifest model.ClusterManifest
	if err := db.DB.Where("version = ?", v.GitVersion).Order("created_at ASC").First(&manifest).Error; err != nil {
		return clusterInfo, err
	}
	clusterInfo.Version = manifest.Name

	// load Node
	kubeNodes, err := kubeClient.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return clusterInfo, err
	}
	for i, node := range kubeNodes.Items {
		if i == 0 {
			if strings.Contains(node.Status.NodeInfo.ContainerRuntimeVersion, "docker") {
				clusterInfo.RuntimeType = "docker"
			} else {
				clusterInfo.RuntimeType = "containerd"
			}
		}
		var item dto.NodeLoadInfo
		if _, ok := node.ObjectMeta.Labels["node-role.kubernetes.io/master"]; ok {
			item.Role = "master"
		} else {
			item.Role = "worker"
		}
		if strings.Contains(node.ObjectMeta.Name, "-") {
			item.Name = strings.Replace(node.ObjectMeta.Name, strings.Split(node.ObjectMeta.Name, "-")[0], loadInfo.Name, 1)
		} else {
			item.Name = node.ObjectMeta.Name
		}
		item.Architecture = node.Status.NodeInfo.Architecture
		item.Port = 22
		for _, addr := range node.Status.Addresses {
			if addr.Type == "InternalIP" {
				item.Ip = addr.Address
			}
		}
		clusterInfo.Nodes = append(clusterInfo.Nodes, item)
	}

	// load kubeadm-config
	kubeAdmMap, err := kubeClient.CoreV1().ConfigMaps("kube-system").Get(context.TODO(), "kubeadm-config", metav1.GetOptions{})
	if err != nil {
		return clusterInfo, fmt.Errorf("can not load kubeadm-config from cluster: %s", err.Error())
	}
	admConfig := kubeAdmMap.Data["ClusterConfiguration"]
	admCfy := make(map[interface{}]interface{})
	if err := yaml.Unmarshal([]byte(admConfig), &admCfy); err != nil {
		return clusterInfo, fmt.Errorf("kubeadm-config yaml unmarshall failed: %s", err.Error())
	}
	admInterface := dyno.ConvertMapI2MapS(admCfy)
	kk, err := json.Marshal(admInterface)
	if err != nil {
		return clusterInfo, fmt.Errorf("kubeadm-config json marshall failed: %s", err.Error())
	}
	var data admConfigStruct
	if err := json.Unmarshal(kk, &data); err != nil {
		return clusterInfo, fmt.Errorf("kubeadm-config json unmarshall failed: %s", err.Error())
	}
	clusterInfo.KubeServiceNodePortRange = data.ApiServer.ExtraArgs.ServiceNodePortRange
	mask, _ := strconv.Atoi(data.Controller.ExtraArgs.NodeCidrMaskSize)
	clusterInfo.KubeNetworkNodePrefix = mask
	clusterInfo.KubeMaxPods = maxNodePodNumMap[mask]
	clusterInfo.KubePodSubnet = data.Network.PodSubnet
	clusterInfo.MaxNodePodNum = 2 << (31 - mask)
	if strings.Contains(clusterInfo.KubePodSubnet, "/") {
		subnets := strings.Split(clusterInfo.KubePodSubnet, "/")
		podMask, _ := strconv.Atoi(subnets[1])
		clusterInfo.MaxNodeNum = (2 << (31 - podMask)) / clusterInfo.MaxNodePodNum
	}

	clusterInfo.KubePodSubnet = data.Network.PodSubnet
	clusterInfo.KubeServiceSubnet = data.Network.ServiceSubnet

	// load kube-proxy
	kubeProxyMap, err := kubeClient.CoreV1().ConfigMaps("kube-system").Get(context.TODO(), "kube-proxy", metav1.GetOptions{})
	if err != nil {
		return clusterInfo, fmt.Errorf("can not load kube-proxy from cluster: %s", err.Error())
	}
	proxyConfig := kubeProxyMap.Data["config.conf"]
	proxyCfy := make(map[interface{}]interface{})
	if err := yaml.Unmarshal([]byte(proxyConfig), &proxyCfy); err != nil {
		return clusterInfo, fmt.Errorf("kube-proxy yaml unmarshall failed: %s", err.Error())
	}
	proxyInterface := dyno.ConvertMapI2MapS(proxyCfy)
	kk2, err := json.Marshal(proxyInterface)
	if err != nil {
		return clusterInfo, fmt.Errorf("kube-proxy json marshall failed: %s", err.Error())
	}
	var data2 proxyConfigStruct
	if err := json.Unmarshal(kk2, &data2); err != nil {
		return clusterInfo, fmt.Errorf("kube-proxy json unmarshall failed: %s", err.Error())
	}
	clusterInfo.KubeProxyMode = data2.Mode
	clusterInfo.NodeportAddress = data2.NodePortAddresses

	// load network
	clusterInfo.EnableDnsCache = "disable"
	daemonsets, err := kubeClient.AppsV1().DaemonSets("kube-system").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return clusterInfo, fmt.Errorf("load daemonsets failed: %s", err.Error())
	}
	for _, daemonset := range daemonsets.Items {
		if daemonset.ObjectMeta.Name == "calico-node" {
			clusterInfo.NetworkType = "calico"
		}
		if daemonset.ObjectMeta.Name == "kube-flannel-ds" {
			clusterInfo.NetworkType = "flannel"
		}
		if daemonset.ObjectMeta.Name == "cilium" {
			clusterInfo.NetworkType = "cilium"
		}
		if daemonset.ObjectMeta.Name == "node-local-dns" {
			clusterInfo.EnableDnsCache = "enable"
		}
		if daemonset.ObjectMeta.Name == "nginx-ingress-controller" {
			clusterInfo.IngressControllerType = "nginx"
		}
		if daemonset.ObjectMeta.Name == "traefik" {
			clusterInfo.IngressControllerType = "traefik"
		}
	}
	return clusterInfo, nil
}

type admConfigStruct struct {
	ApiServer  apiServerStruct  `json:"apiServer"`
	Controller ControllerStruct `json:"controllerManager"`
	Network    networkStruct    `json:"networking"`
}

type proxyConfigStruct struct {
	Mode              string `json:"mode"`
	NodePortAddresses string `json:"nodePortAddresses"`
}

type apiServerStruct struct {
	ExtraArgs extraArgsStruct `json:"extraArgs"`
}

type ControllerStruct struct {
	ExtraArgs extraArgsStruct `json:"extraArgs"`
}

type networkStruct struct {
	PodSubnet     string `json:"podSubnet"`
	ServiceSubnet string `json:"serviceSubnet"`
}

type extraArgsStruct struct {
	ServiceNodePortRange string `json:"service-node-port-range"`
	NodeCidrMaskSize     string `json:"node-cidr-mask-size"`
}
