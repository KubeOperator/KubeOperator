package service

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/logger"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/repository"
	clusterUtil "github.com/KubeOperator/KubeOperator/pkg/util/cluster"
	"github.com/KubeOperator/KubeOperator/pkg/util/ipaddr"
	"github.com/KubeOperator/KubeOperator/pkg/util/ssh"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var (
	CheckHostSSHConnection = "CHECK_HOST_SSH_CONNECTION"
	CheckK8sToken          = "CHECK_K8S_TOKEN"
	CheckK8sAPI            = "CHECK_K8S_API"
	CheckK8sNodeStatus     = "CHECK_K8S_NODE_STATUS"
	CheckKubeRouter        = "CHECK_KUBE_ROUTER"

	StatusSuccess        = "STATUS_SUCCESS"
	StatusWarning        = "STATUS_WARNING"
	StatusFailed         = "STATUS_FAILED"
	StatusError          = "STATUS_ERROR"
	StatusSolvedManually = "STATUS_SOLVED_MANUALLY"
	StatusRecoverd       = "STATUS_RECOVERD"

	RecoverNodeStatus   = "RECOVER_SYNC_NODE_STATUS"
	RecoverSyncRouterIP = "RECOVER_SYNC_ROUTER_IP"
	RecoverSyncToken    = "RECOVER_SYNC_TOKEN"
	RecoverConnConf     = "RECOVER_CONN_CONF"
	RecoverHostConn     = "RECOVER_HOST_CONN"
	RecoverAPIConn      = "RECOVER_API_CONN"
)

type ClusterHealthService interface {
	HealthCheck(clusterName string) (*dto.ClusterHealth, error)
	Recover(clusterName string, ch dto.ClusterHealth) ([]dto.ClusterRecoverItem, error)
}

type clusterHealthService struct {
	clusterService     ClusterService
	clusterRepo        repository.ClusterRepository
	clusterNodeRepo    repository.ClusterNodeRepository
	clusterInitService ClusterInitService
}

func NewClusterHealthService() ClusterHealthService {
	return &clusterHealthService{
		clusterService:     NewClusterService(),
		clusterRepo:        repository.NewClusterRepository(),
		clusterNodeRepo:    repository.NewClusterNodeRepository(),
		clusterInitService: NewClusterInitService(),
	}
}

type HealthCheckFunc func(c model.Cluster) dto.ClusterHealthHook

func (c clusterHealthService) HealthCheck(clusterName string) (*dto.ClusterHealth, error) {
	clu, err := c.clusterRepo.GetWithPreload(clusterName, []string{"SpecConf", "SpecNetwork", "SpecRuntime", "Secret", "Nodes", "Nodes.Host", "Nodes.Host.Credential"})
	if err != nil {
		return nil, err
	}
	results := dto.ClusterHealth{Level: StatusError}
	results.Level = StatusError
	if clu.Source != constant.ClusterSourceExternal {
		sshResult := checkHostSSHConnected(clu)
		results.Hooks = append(results.Hooks, sshResult)
		if sshResult.Level == StatusError {
			return &results, nil
		}
	}
	client, tokenResult := checkClusterConn(clu)
	if tokenResult.Level == StatusError {
		results.Hooks = append(results.Hooks, tokenResult)
		return &results, nil
	}
	results.Hooks = append(results.Hooks, tokenResult)

	apiResult := checkKubernetesApi(clu)
	results.Hooks = append(results.Hooks, apiResult)
	if apiResult.Level == StatusError {
		return &results, nil
	}

	nodes, nodeResult := checkKubernetesNodeStatus(clu, client)
	if nodeResult.Level == StatusError {
		for _, node := range nodes {
			for _, addr := range node.Status.Addresses {
				if addr.Type == "InternalIP" {
					nodeResult.AdjustValue += nodeResult.AdjustValue + addr.Address + ","
				}
			}
		}
		results.Hooks = append(results.Hooks, nodeResult)
		return &results, nil
	}
	results.Hooks = append(results.Hooks, nodeResult)

	routerResult := checkKubeRouter(clu, nodes)
	if routerResult.Level == StatusError {
		isExist := false
		for _, node := range nodes {
			if _, ok := node.ObjectMeta.Labels["node-role.kubernetes.io/master"]; !ok {
				continue
			}
			for _, addr := range node.Status.Addresses {
				if addr.Type == "InternalIP" {
					routerResult.AdjustValue = addr.Address
					isExist = true
					break
				}
			}
			if isExist {
				break
			}
		}
		results.Hooks = append(results.Hooks, routerResult)
		return &results, nil
	}
	results.Hooks = append(results.Hooks, routerResult)

	results.Level = StatusSuccess
	return &results, nil
}

// 检查各个主机连接状态，不存在可用主节点时错误
func checkHostSSHConnected(c model.Cluster) dto.ClusterHealthHook {
	result := dto.ClusterHealthHook{
		Name:  CheckHostSSHConnection,
		Level: StatusSuccess,
	}
	aliveMaster := 0
	wg := sync.WaitGroup{}
	for i := range c.Nodes {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			if err := ipaddr.Ping(c.Nodes[n].Host.Ip); err != nil {
				result.Level = StatusWarning
				result.Msg += fmt.Sprintf("Ping %s failed: %s,", c.Nodes[n].Host.Ip, err.Error())
				return
			}
			sshCfg := c.Nodes[n].ToSSHConfig()
			sshClient, err := ssh.New(&sshCfg)
			if err != nil {
				result.Level = StatusWarning
				result.Msg += fmt.Sprintf("SSH %s failed: %s,", c.Nodes[n].Host.Ip, err.Error())
				return
			}
			if err := sshClient.Ping(); err != nil {
				result.Level = StatusWarning
				result.Msg += fmt.Sprintf("SSH ping %s failed: %s,", c.Nodes[n].Host.Ip, err.Error())
				return
			}
			if c.Nodes[n].Role == constant.NodeRoleNameMaster {
				aliveMaster++
			}
		}(i)
	}
	wg.Wait()
	if aliveMaster == 0 {
		result.Level = StatusError
		result.Msg = "no available master in cluster"
	}
	return result
}

// 检查集群可连接性 (token / config / cret)
func checkClusterConn(c model.Cluster) (*kubernetes.Clientset, dto.ClusterHealthHook) {
	result := dto.ClusterHealthHook{
		Name:  CheckK8sToken,
		Level: StatusSuccess,
	}
	client, err := clusterUtil.NewClusterClient(&c)
	if err != nil {
		result.Msg = fmt.Sprintf("Get token form cluster failed %s", err.Error())
		result.Level = StatusError
		return client, result
	}
	return client, result
}

// 用 lb_ip 去请求集群 healthz 接口，判断 api 可用性
func checkKubernetesApi(c model.Cluster) dto.ClusterHealthHook {
	result := dto.ClusterHealthHook{
		Name:  CheckK8sAPI,
		Level: StatusSuccess,
	}
	isOK, err := GetClusterStatusByAPI(fmt.Sprintf("%s:%d", c.SpecConf.LbKubeApiserverIp, c.SpecConf.KubeApiServerPort), &c)
	if !isOK {
		result.Msg = err
		result.Level = StatusError
	}
	return result
}

// 检查集群节点数量与数据库节点数量
func checkKubernetesNodeStatus(c model.Cluster, client *kubernetes.Clientset) ([]v1.Node, dto.ClusterHealthHook) {
	var nodes []model.ClusterNode
	result := dto.ClusterHealthHook{
		Name:  CheckK8sNodeStatus,
		Level: StatusSuccess,
	}

	kubeNodes, err := client.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Log.Errorf("get cluster %s kubeNodes error %s", c.Name, err.Error())
		result.Msg = fmt.Sprintf("get cluster %s kubeNodes error %s", c.Name, err.Error())
		result.Level = StatusError
		return nil, result
	}
	if err := db.DB.Where("cluster_id = ?", c.ID).Find(&nodes).Error; err != nil {
		logger.Log.Errorf("get cluster %s nodes from db error %s", c.Name, err.Error())
		result.Msg = fmt.Sprintf("get cluster %s nodes from db error %s", c.Name, err.Error())
		result.Level = StatusError
		return nil, result
	}
	if len(nodes) != len(kubeNodes.Items) {
		logger.Log.Errorf("The number of system nodes: %d does not match the number of k8s nodes: %d", len(nodes), len(kubeNodes.Items))
		result.Msg = fmt.Sprintf("The number of system nodes: %d does not match the number of k8s nodes: %d", len(nodes), len(kubeNodes.Items))
		result.Level = StatusError
		return nil, result
	}

	return kubeNodes.Items, result
}

// 检查 kuberouter 连接
func checkKubeRouter(c model.Cluster, nodes []v1.Node) dto.ClusterHealthHook {
	result := dto.ClusterHealthHook{
		Name:  CheckKubeRouter,
		Level: StatusSuccess,
	}
	isExist := false
	for _, node := range nodes {
		for _, addr := range node.Status.Addresses {
			if addr.Type == "InternalIP" {
				if addr.Address == c.SpecConf.KubeRouter {
					isExist = true
				}
			}
		}
	}
	if !isExist {
		result.Msg = fmt.Sprintf("The address %s of kube router is not alived in cluster", c.SpecConf.KubeRouter)
		result.Level = StatusError
	}
	return result
}

func (c clusterHealthService) Recover(clusterName string, ch dto.ClusterHealth) ([]dto.ClusterRecoverItem, error) {
	var result []dto.ClusterRecoverItem
	clu, err := c.clusterRepo.GetWithPreload(clusterName, []string{"SpecConf", "SpecNetwork", "SpecRuntime", "Secret", "Nodes", "Nodes.Host", "Nodes.Host.Credential"})
	if err != nil {
		return result, err
	}

	for i := range ch.Hooks {
		if ch.Hooks[i].Level == StatusError {
			ri := dto.ClusterRecoverItem{
				Name: ch.Hooks[i].Name,
			}
			switch ch.Hooks[i].Name {
			case CheckHostSSHConnection:
				ri.Result = StatusSolvedManually
				ri.Method = RecoverHostConn
				result = append(result, ri)
				return result, nil
			case CheckK8sAPI:
				c.recoverK8sAPI(clu, &ri)
				result = append(result, ri)
				return result, nil
			case CheckK8sToken:
				if clu.SpecConf.AuthenticationMode == constant.AuthenticationModeBearer {
					ri.Method = RecoverSyncToken
					if err := c.clusterInitService.GatherKubernetesToken(clu); err != nil {
						ri.Result = StatusFailed
						ri.Msg = err.Error()
						result = append(result, ri)
						return result, nil
					} else {
						ri.Result = StatusRecoverd
						result = append(result, ri)
						return result, nil
					}
				} else {
					ri.Result = StatusSolvedManually
					ri.Method = RecoverConnConf
					result = append(result, ri)
					return result, nil
				}
			case CheckK8sNodeStatus:
				c.recoverNodeStatus(clu, &ri, ch.Hooks[i].AdjustValue)
				result = append(result, ri)
				return result, nil
			case CheckKubeRouter:
				c.recoverKubeRouter(clu, &ri, ch.Hooks[i].AdjustValue)
				result = append(result, ri)
				return result, nil
			default:
				return result, nil
			}
		}
	}

	return result, nil
}

// 主节点中筛选一个存活的主机，修改为 lb_kube_apiserver_ip
// vip 时不操作
func (c clusterHealthService) recoverK8sAPI(m model.Cluster, ri *dto.ClusterRecoverItem) {
	var hosts []string
	ri.Method = RecoverAPIConn
	if m.SpecConf.LbMode == constant.ClusterSourceExternal || m.Source == constant.ClusterSourceExternal {
		ri.Result = StatusSolvedManually
		return
	}
	port := m.SpecConf.KubeApiServerPort
	masters, err := c.clusterNodeRepo.AllMaster(m.ID)
	if err != nil {
		ri.Result = StatusFailed
		ri.Msg = fmt.Sprintf("get master error %s", err.Error())
		return
	}
	for i := range masters {
		hosts = append(hosts, fmt.Sprintf("%s:%d", masters[i].Host.Ip, port))
	}

	aliveHost, err := clusterUtil.SelectAliveHost(hosts)
	if err != nil {
		ri.Result = StatusFailed
		ri.Msg = fmt.Sprintf("select alive host error %s", err.Error())
		return
	}
	isOk, msg := GetClusterStatusByAPI(string(aliveHost), &m)
	if isOk {
		if err := db.DB.Model(&model.ClusterSpecConf{}).Where("cluster_id = ?", m.ID).Updates(map[string]interface{}{"lb_kube_apiserver_ip": strings.Split(string(aliveHost), ":")[0]}).Error; err != nil {
			ri.Result = StatusFailed
			ri.Msg = err.Error()
			return
		}
		ri.Method = RecoverSyncRouterIP
		ri.Result = StatusRecoverd
		return
	}
	ri.Result = StatusFailed
	ri.Msg = msg
}

// 主节点中筛选一个存活的主机，修改为 kube_router
func (c clusterHealthService) recoverKubeRouter(m model.Cluster, ri *dto.ClusterRecoverItem, adjustValue string) {
	ri.Method = RecoverSyncRouterIP
	kubeRouter := ""
	if len(adjustValue) != 0 {
		kubeRouter = adjustValue
	} else {
		client, err := clusterUtil.NewClusterClient(&m)
		if err != nil {
			ri.Result = StatusFailed
			ri.Msg = err.Error()
			return
		}
		kubeNodes, err := client.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			ri.Result = StatusFailed
			ri.Msg = err.Error()
			return
		}
		isExist := false
		for _, node := range kubeNodes.Items {
			if _, ok := node.ObjectMeta.Labels["node-role.kubernetes.io/master"]; !ok {
				continue
			}
			for _, addr := range node.Status.Addresses {
				if addr.Type == "InternalIP" {
					kubeRouter = addr.Address
					isExist = true
					break
				}
			}
			if isExist {
				break
			}
		}
		if kubeRouter == "" {
			ri.Result = StatusFailed
			ri.Msg = "No master available in cluster"
			return
		}
	}
	if err := db.DB.Model(&model.ClusterSpecConf{}).Where("cluster_id = ?", m.ID).Updates(map[string]interface{}{"kube_router": kubeRouter}).Error; err != nil {
		ri.Result = StatusFailed
		ri.Msg = err.Error()
		return
	}
	ri.Result = StatusRecoverd
}

// 节点数量同步，将数据库中多出的节点标记为脏数据 且修改为失联状态
func (c clusterHealthService) recoverNodeStatus(m model.Cluster, ri *dto.ClusterRecoverItem, adjustValue string) {
	ri.Method = RecoverNodeStatus
	var nodes []model.ClusterNode
	if err := db.DB.Where("cluster_id = ?", m.ID).Preload("Host").Find(&nodes).Error; err != nil {
		ri.Result = StatusFailed
		ri.Msg = err.Error()
		return
	}
	var nodeIDs []string
	alivedIP := strings.Split(adjustValue, ",")
	if len(adjustValue) != 0 {
		for _, node := range nodes {
			for _, ip := range alivedIP {
				if ip != "" && ip == node.Host.Ip {
					continue
				}
			}
			nodeIDs = append(nodeIDs, node.ID)
		}
	} else {
		client, err := clusterUtil.NewClusterClient(&m)
		if err != nil {
			ri.Result = StatusFailed
			ri.Msg = err.Error()
			return
		}

		kubeNodes, err := client.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			ri.Result = StatusFailed
			ri.Msg = err.Error()
			return
		}
		for _, node := range nodes {
			hasNode := false
			for _, kn := range kubeNodes.Items {
				if kn.ObjectMeta.Name == node.Name {
					hasNode = true
					break
				}
			}
			if hasNode {
				continue
			}
			nodeIDs = append(nodeIDs, node.ID)
		}
	}
	if err := db.DB.Model(&model.ClusterNode{}).Where("id in (?)", nodeIDs).Updates(map[string]interface{}{"status": constant.StatusLost, "dirty": true}).Error; err != nil {
		ri.Result = StatusFailed
		ri.Msg = err.Error()
		return
	}
	ri.Result = StatusRecoverd
}

func GetClusterStatusByAPI(addr string, cluster *model.Cluster) (bool, string) {
	reqURL := fmt.Sprintf("https://%s/livez", addr)
	conf, err := clusterUtil.LoadConnConf(cluster, addr)
	if err != nil {
		_, _ = false, fmt.Sprintf("load cluster config error %s", err.Error())
	}
	tls2, err := rest.TransportFor(conf)
	if err != nil {
		_, _ = false, fmt.Sprintf("config transport error %s", err.Error())
	}
	client := &http.Client{Timeout: 1 * time.Second, Transport: tls2}
	request, _ := http.NewRequest("GET", reqURL, nil)
	response, err := client.Do(request)
	if err != nil {
		return false, fmt.Sprintf("Https get error %s", err.Error())
	}
	defer response.Body.Close()
	if response.StatusCode == 200 {
		return true, ""
	}
	s, _ := ioutil.ReadAll(response.Body)
	return false, fmt.Sprintf("Api check error %s", string(s))
}
