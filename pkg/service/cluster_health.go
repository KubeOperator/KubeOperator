package service

import (
	"context"
	"crypto/tls"
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
	clusterUtil "github.com/KubeOperator/KubeOperator/pkg/util/cluster"
	"github.com/KubeOperator/KubeOperator/pkg/util/ipaddr"
	kubeUtil "github.com/KubeOperator/KubeOperator/pkg/util/kubernetes"
	"github.com/KubeOperator/KubeOperator/pkg/util/ssh"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
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
)

type ClusterHealthService interface {
	HealthCheck(clusterName string) (*dto.ClusterHealth, error)
	Recover(clusterName string, ch dto.ClusterHealth) ([]dto.ClusterRecoverItem, error)
}

type clusterHealthService struct {
	clusterService     ClusterService
	clusterInitService ClusterInitService
}

func NewClusterHealthService() ClusterHealthService {
	return &clusterHealthService{
		clusterService:     NewClusterService(),
		clusterInitService: NewClusterInitService(),
	}
}

type HealthCheckFunc func(c model.Cluster) dto.ClusterHealthHook

func (c clusterHealthService) HealthCheck(clusterName string) (*dto.ClusterHealth, error) {
	clu, err := c.clusterService.Get(clusterName)
	if err != nil {
		return nil, err
	}
	results := dto.ClusterHealth{Level: StatusError}
	results.Level = StatusError
	if clu.Source == constant.ClusterSourceLocal {
		sshclient, sshResult := checkHostSSHConnected(clu.Cluster)
		results.Hooks = append(results.Hooks, sshResult)
		if sshResult.Level == StatusError {
			return &results, nil
		}

		token, tokenResult := checkKubernetesToken(clu.Cluster, sshclient)
		if tokenResult.Level == StatusError {
			tokenResult.AdjustValue = token
			results.Hooks = append(results.Hooks, tokenResult)
			return &results, nil
		}
		results.Hooks = append(results.Hooks, tokenResult)
	}

	apiResult := checkKubernetesApi(clu.Cluster)
	results.Hooks = append(results.Hooks, apiResult)
	if apiResult.Level == StatusError {
		return &results, nil
	}

	nodes, nodeResult := checkKubernetesNodeStatus(clu.Cluster)
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

	routerResult := checkKubeRouter(clu.Cluster, nodes)
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

func checkHostSSHConnected(c model.Cluster) (*ssh.SSH, dto.ClusterHealthHook) {
	result := dto.ClusterHealthHook{
		Name:  CheckHostSSHConnection,
		Level: StatusSuccess,
	}
	var backSSHClient *ssh.SSH
	isExist := false
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
				if !isExist {
					backSSHClient = sshClient
					isExist = true
				}
				aliveMaster++
			}
		}(i)
	}
	wg.Wait()
	if !(aliveMaster > 0) {
		result.Level = StatusError
	}
	return backSSHClient, result
}

func checkKubernetesToken(c model.Cluster, sshClient *ssh.SSH) (string, dto.ClusterHealthHook) {
	clusterService := NewClusterService()
	result := dto.ClusterHealthHook{
		Name:  CheckK8sToken,
		Level: StatusSuccess,
	}
	token, err := clusterUtil.GetClusterTokenWithoutRetry(sshClient)
	if err != nil {
		result.Msg = fmt.Sprintf("Get token form cluster failed %s", err.Error())
		result.Level = StatusError
		return "", result
	}
	secret, err := clusterService.GetSecrets(c.Name)
	if err != nil {
		result.Msg = fmt.Sprintf("Get token from db failed %s", err.Error())
		result.Level = StatusError
		return token, result
	}
	if token != secret.KubernetesToken {
		result.Msg = "The cluster token is inconsistent with the database"
		result.Level = StatusError
		return token, result
	}
	return token, result
}

func checkKubernetesApi(c model.Cluster) dto.ClusterHealthHook {
	result := dto.ClusterHealthHook{
		Name:  CheckK8sAPI,
		Level: StatusSuccess,
	}
	isOK, err := GetClusterStatusByAPI(c)
	if !isOK {
		result.Msg = err
		result.Level = StatusError
	}
	return result
}

func checkKubernetesNodeStatus(c model.Cluster) ([]v1.Node, dto.ClusterHealthHook) {
	var nodes []model.ClusterNode
	client, level, msg := getBaseParams(c)
	result := dto.ClusterHealthHook{
		Name:  CheckK8sNodeStatus,
		Level: level,
		Msg:   msg,
	}
	if len(msg) != 0 {
		logger.Log.Errorf("get cluster %s base info failed: %s", c.Name, msg)
		return nil, result
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

func checkKubeRouter(c model.Cluster, nodes []v1.Node) dto.ClusterHealthHook {
	result := dto.ClusterHealthHook{
		Name:  CheckKubeRouter,
		Level: StatusSuccess,
	}
	isExist := false
	for _, node := range nodes {
		for _, addr := range node.Status.Addresses {
			if addr.Type == "InternalIP" {
				if addr.Address == c.Spec.KubeRouter {
					isExist = true
				}
			}
		}
	}
	if !isExist {
		result.Msg = fmt.Sprintf("The address %s of kube router is not alived in cluster", c.Spec.KubeRouter)
		result.Level = StatusError
	}
	return result
}

var resolveMethods = map[string]string{
	CheckHostSSHConnection: "CHECK_SSH_CONNECTION",
	CheckK8sToken:          "GET_K8S_TOKEN_ANGIN",
	CheckK8sAPI:            "CHECK_API_CONNECTION",
	CheckK8sNodeStatus:     "UPDATE_CLUSTER_NODE_STATUS",
	CheckKubeRouter:        "UPDATE_KUBE_ROUTER",
}

func (c clusterHealthService) Recover(clusterName string, ch dto.ClusterHealth) ([]dto.ClusterRecoverItem, error) {
	var result []dto.ClusterRecoverItem
	clu, err := c.clusterService.Get(clusterName)
	if err != nil {
		return result, err
	}
	switch ch.Level {
	case StatusError:
		for i := range ch.Hooks {
			if ch.Hooks[i].Level == StatusError {
				ri := dto.ClusterRecoverItem{
					Name:   ch.Hooks[i].Name,
					Method: resolveMethods[ch.Hooks[i].Name],
				}
				switch ch.Hooks[i].Name {
				case CheckHostSSHConnection, CheckK8sAPI:
					ri.Result = StatusSolvedManually
					result = append(result, ri)
					return result, nil
				case CheckK8sToken:
					if len(ch.Hooks[i].AdjustValue) != 0 {
						if err := db.DB.Model(&model.ClusterSecret{}).Where("id = ?", clu.SecretID).Updates(map[string]interface{}{"kubernetes_token": ch.Hooks[i].AdjustValue}).Error; err != nil {
							ri.Result = StatusFailed
							ri.Msg = err.Error()
							result = append(result, ri)
							return result, nil
						}
					} else {
						if err := c.clusterInitService.GatherKubernetesToken(clu.Cluster); err != nil {
							ri.Result = StatusFailed
							ri.Msg = err.Error()
							result = append(result, ri)
							return result, nil
						}
					}
					ri.Result = StatusRecoverd
					result = append(result, ri)
				case CheckK8sNodeStatus:
					var nodes []model.ClusterNode
					if err := db.DB.Where("cluster_id = ?", clu.Cluster.ID).Preload("Host").Find(&nodes).Error; err != nil {
						ri.Result = StatusFailed
						ri.Msg = err.Error()
						result = append(result, ri)
						return result, nil
					}
					var nodeIDs []string
					alivedIP := strings.Split(ch.Hooks[i].AdjustValue, ",")
					if len(ch.Hooks[i].AdjustValue) != 0 {
						for _, node := range nodes {
							for _, ip := range alivedIP {
								if ip != "" && ip == node.Host.Ip {
									continue
								}
							}
							nodeIDs = append(nodeIDs, node.ID)
						}
					} else {
						client, _, msg := getBaseParams(clu.Cluster)
						if len(msg) != 0 {
							ri.Result = StatusFailed
							ri.Msg = err.Error()
							result = append(result, ri)
							return result, nil
						}

						kubeNodes, err := client.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
						if err != nil {
							ri.Result = StatusFailed
							ri.Msg = err.Error()
							result = append(result, ri)
							return result, nil
						}
						for _, node := range nodes {
							hasNode := false
							for _, kn := range kubeNodes.Items {
								for _, addr := range kn.Status.Addresses {
									if addr.Type == "InternalIP" && node.Host.Ip == addr.Address {
										hasNode = true
										break
									}
								}
								if hasNode {
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
						result = append(result, ri)
						return result, nil
					}
					ri.Result = StatusSolvedManually
					result = append(result, ri)
				case CheckKubeRouter:
					kubeRouter := ""
					if len(ch.Hooks[i].AdjustValue) != 0 {
						kubeRouter = ch.Hooks[i].AdjustValue
					} else {
						client, _, msg := getBaseParams(clu.Cluster)
						if len(msg) != 0 {
							ri.Result = StatusFailed
							ri.Msg = err.Error()
							result = append(result, ri)
							return result, nil
						}
						kubeNodes, err := client.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
						if err != nil {
							ri.Result = StatusFailed
							ri.Msg = err.Error()
							result = append(result, ri)
							return result, nil
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
					}
					if kubeRouter == "" {
						ri.Result = StatusFailed
						ri.Msg = "No address available in cluster"
						result = append(result, ri)
						return result, nil
					}
					if err := db.DB.Model(&model.ClusterSpec{}).Where("id = ?", clu.SpecID).Updates(map[string]interface{}{"kube_router": kubeRouter}).Error; err != nil {
						ri.Result = StatusFailed
						ri.Msg = err.Error()
						result = append(result, ri)
						return result, nil
					}
					ri.Result = StatusRecoverd
					result = append(result, ri)
				default:
					return result, nil
				}
			}
		}
	}

	return result, nil
}

func getBaseParams(c model.Cluster) (*kubernetes.Clientset, string, string) {
	clusterService := NewClusterService()
	secret, err := clusterService.GetSecrets(c.Name)
	if err != nil {
		msg := fmt.Sprintf("get cluster %s secret error %s", c.Name, err.Error())
		level := StatusError
		return nil, level, msg
	}

	endpoints, err := clusterService.GetApiServerEndpoints(c.Name)
	if err != nil {
		msg := fmt.Sprintf("get cluster %s endpoint error %s", c.Name, err.Error())
		level := StatusError
		return nil, level, msg
	}

	_, err = kubeUtil.SelectAliveHost(endpoints)
	if err != nil {
		msg := fmt.Sprintf("no alived host in cluster %s", c.Name)
		level := StatusError
		return nil, level, msg
	}

	kubeClient, err := kubeUtil.NewKubernetesClient(&kubeUtil.Config{
		Hosts: endpoints,
		Token: secret.KubernetesToken,
	})
	if err != nil {
		msg := fmt.Sprintf("get cluster %s kubeclient error %s", c.Name, err.Error())
		level := StatusError
		return nil, level, msg
	}

	return kubeClient, StatusSuccess, ""
}

func GetClusterStatusByAPI(c model.Cluster) (bool, string) {
	clusterService := NewClusterService()
	endpoints, err := clusterService.GetApiServerEndpoints(c.Name)
	if err != nil {
		return false, fmt.Sprintf("Get cluster secret error %s", err.Error())
	}
	aliveHost, err := kubeUtil.SelectAliveHost(endpoints)
	if err != nil {
		return false, fmt.Sprintf("Select alive host error %s", err.Error())
	}
	reqURL := fmt.Sprintf("https://%s/livez", aliveHost)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Timeout: 1 * time.Second, Transport: tr}
	secret, err := clusterService.GetSecrets(c.Name)
	if err != nil {
		return false, fmt.Sprintf("Get secrets error %s", err.Error())
	}
	token := fmt.Sprintf("%s %s", "Bearer", secret.KubernetesToken)
	request, _ := http.NewRequest("GET", reqURL, nil)
	request.Header.Add("Authorization", token)
	response, err := client.Do(request)
	if err != nil {
		return false, fmt.Sprintf("Http get error %s", err.Error())
	}
	defer response.Body.Close()
	if response.StatusCode == 200 {
		return true, ""
	}
	s, _ := ioutil.ReadAll(response.Body)
	return false, fmt.Sprintf("Api check error %s", string(s))
}
