package service

import (
	"context"
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/repository"
	kubeUtil "github.com/KubeOperator/KubeOperator/pkg/util/kubernetes"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"strconv"
	"strings"
	"sync"
)

type ClusterImportService interface {
	Import(clusterImport dto.ClusterImport) error
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
	var address string
	var port int
	if strings.HasSuffix(clusterImport.ApiServer, "/") {
		clusterImport.ApiServer = strings.Replace(clusterImport.ApiServer, "/", "", -1)
	}
	if strings.Contains(clusterImport.ApiServer, "http://") {
		clusterImport.ApiServer = strings.Replace(clusterImport.ApiServer, "http://", "", -1)
	}
	if strings.Contains(clusterImport.ApiServer, "https://") {
		clusterImport.ApiServer = strings.Replace(clusterImport.ApiServer, "https://", "", -1)
	}
	if strings.Contains(clusterImport.ApiServer, ":") {
		strs := strings.Split(clusterImport.ApiServer, ":")
		address = strs[0]
		port, _ = strconv.Atoi(strs[1])
	} else {
		address = clusterImport.ApiServer
		port = 80
	}
	cluster := model.Cluster{
		Name:   clusterImport.Name,
		Source: constant.ClusterSourceExternal,
		Status: model.ClusterStatus{
			Phase: constant.ClusterRunning,
		},
		Spec: model.ClusterSpec{
			LbKubeApiserverIp: address,
			KubeApiServerPort: port,
			KubeRouter:        clusterImport.Router,
		},
		Secret: model.ClusterSecret{
			KubeadmToken:    "",
			KubernetesToken: clusterImport.Token,
		},
	}
	if err := gatherClusterInfo(&cluster); err != nil {
		return err
	}
	if err := c.clusterRepo.Save(&cluster); err != nil {
		return err
	}
	project, err := c.projectRepository.Get(clusterImport.ProjectName)
	if err != nil {
		return err

	}
	if err := c.projectResourceRepository.Create(model.ProjectResource{
		ResourceID:   cluster.ID,
		ProjectID:    project.ID,
		ResourceType: constant.ResourceCluster,
	}); err != nil {
		return err
	}
	return nil
}

func gatherClusterInfo(cluster *model.Cluster) error {
	c, err := kubeUtil.NewKubernetesClient(&kubeUtil.Config{
		Hosts: []kubeUtil.Host{kubeUtil.Host(cluster.Spec.LbKubeApiserverIp)},
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
		log.Error(err.Error())
		return
	}
	cluster.Spec.Version = v.GitVersion
}

func getKubeNodes(cluster *model.Cluster, client *kubernetes.Clientset, wg *sync.WaitGroup) {
	defer wg.Done()
	nodes, err := client.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Error(err.Error())
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
		log.Error(err.Error())
		return
	}
	networkMap := map[string]int{
		"flannel": 0,
		"calico":  0,
	}
	for _, dp := range dps.Items {
		for i := range networkMap {
			if strings.Contains(dp.Name, i) {
				networkMap[i] += 1
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
		log.Error(err.Error())
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
