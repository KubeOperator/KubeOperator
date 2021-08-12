package service

import (
	"context"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var (
	CheckHostSSHConnection = "CHECK_HOST_SSH_CONNECTION"
	CheckK8sToken          = "CHECK_K8S_TOKEN"
	CheckK8sAPI            = "CHECK_K8S_API"
	CheckK8sNodeStatus     = "CHECK_K8S_NODE_STATUS"

	StatusSuccess = "STATUS_SUCCESS"
	StatusWarning = "STATUS_WARNING"
	StatusFailed  = "STATUS_FAILED"
	StatusError   = "STATUS_ERROR"
)

type ClusterHealthService interface {
	HealthCheck(clusterName string) (*dto.ClusterHealth, error)
	Recover(clusterName string) ([]dto.ClusterRecoverItem, error)
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
	hookList := []HealthCheckFunc{
		checkHostSSHConnected,
		checkKubernetesToken,
		checkKubernetesApi,
		checkKubernetesNodeStatus,
	}
	clu, err := c.clusterService.Get(clusterName)
	if err != nil {
		return nil, err
	}
	result := dto.ClusterHealth{}
	for _, h := range hookList {
		hookResult := h(clu.Cluster)
		if hookResult.Level == StatusError {
			result.Level = StatusError
			result.Hooks = append(result.Hooks, hookResult)
			return &result, nil
		}
		result.Hooks = append(result.Hooks, hookResult)
	}
	return &result, nil
}

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
	if !(aliveMaster > 0) {
		result.Level = StatusError
	}
	return result
}

func checkKubernetesToken(c model.Cluster) dto.ClusterHealthHook {
	clusterService := NewClusterService()
	result := dto.ClusterHealthHook{
		Name:  CheckK8sToken,
		Level: StatusSuccess,
	}
	token, err := getClusterToken(c)
	if err != nil {
		result.Msg = fmt.Sprintf("Get token form cluster failed %s", err.Error())
		result.Level = StatusError
		return result
	}
	secret, err := clusterService.GetSecrets(c.Name)
	if err != nil {
		result.Msg = fmt.Sprintf("Get token from db failed %s", err.Error())
		result.Level = StatusError
		return result
	}
	if token != secret.KubernetesToken {
		result.Msg = "The cluster token is inconsistent with the database"
		result.Level = StatusError
		return result
	}
	return result
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

func checkKubernetesNodeStatus(c model.Cluster) dto.ClusterHealthHook {
	var nodes []model.ClusterNode
	client, level, msg := getBaseParams(c)
	result := dto.ClusterHealthHook{
		Name:  CheckK8sNodeStatus,
		Level: level,
		Msg:   msg,
	}
	if len(msg) != 0 {
		logger.Log.Errorf("get cluster %s base info failed: %s", c.Name, msg)
		return result
	}

	kubeNodes, err := client.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Log.Errorf("get cluster %s kubeNodes error %s", c.Name, err.Error())
		result.Msg = fmt.Sprintf("get cluster %s kubeNodes error %s", c.Name, err.Error())
		result.Level = StatusError
		return result
	}
	if err := db.DB.Where("cluster_id = ?", c.ID).Find(&nodes).Error; err != nil {
		logger.Log.Errorf("get cluster %s nodes from db error %s", c.Name, err.Error())
		result.Msg = fmt.Sprintf("get cluster %s nodes from db error %s", c.Name, err.Error())
		result.Level = StatusError
		return result
	}
	if len(nodes) != len(kubeNodes.Items) {
		logger.Log.Errorf("The number of system nodes: %d does not match the number of k8s nodes: %d", len(nodes), len(kubeNodes.Items))
		result.Msg = fmt.Sprintf("The number of system nodes: %d does not match the number of k8s nodes: %d", len(nodes), len(kubeNodes.Items))
		result.Level = StatusError
		return result
	}

	return result
}

var resolveMethods = map[string]string{
	CheckHostSSHConnection: "NO_METHODS",
	CheckK8sToken:          "GET_K8S_TOKEN_ANGIN",
	CheckK8sAPI:            "NO_METHODS",
	CheckK8sNodeStatus:     "UPDATE_CLUSTER_NODE_STATUS",
}

func (c clusterHealthService) Recover(clusterName string) ([]dto.ClusterRecoverItem, error) {
	var result []dto.ClusterRecoverItem
	clu, err := c.clusterService.Get(clusterName)
	if err != nil {
		return result, err
	}
	ch, err := c.HealthCheck(clusterName)
	if err != nil {
		return result, err
	}
	switch ch.Level {
	case StatusError:
		for i := range ch.Hooks {
			if ch.Hooks[i].Level == StatusError {
				switch ch.Hooks[i].Name {
				case CheckHostSSHConnection, CheckK8sAPI:
					ri := dto.ClusterRecoverItem{
						Name:     resolveMethods[ch.Hooks[i].Name],
						HookName: ch.Hooks[i].Name,
						Result:   StatusFailed,
						Msg:      "No method",
					}
					result = append(result, ri)
					return result, nil
				case CheckK8sToken:
					ri := dto.ClusterRecoverItem{
						Name:     resolveMethods[ch.Hooks[i].Name],
						HookName: ch.Hooks[i].Name,
					}
					err := c.clusterInitService.GatherKubernetesToken(clu.Cluster)
					if err != nil {
						ri.Result = StatusFailed
						ri.Msg = err.Error()
						result = append(result, ri)
						return result, nil
					}
					ri.Result = StatusSuccess
					result = append(result, ri)
				case CheckK8sNodeStatus:
					client, _, msg := getBaseParams(clu.Cluster)
					ri := dto.ClusterRecoverItem{
						Name:     resolveMethods[ch.Hooks[i].Name],
						HookName: ch.Hooks[i].Name,
					}
					if len(msg) != 0 {
						ri.Result = StatusFailed
						ri.Msg = err.Error()
						result = append(result, ri)
						return result, nil
					}

					var nodes []model.ClusterNode
					kubeNodes, err := client.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
					if err != nil {
						ri.Result = StatusFailed
						ri.Msg = err.Error()
						result = append(result, ri)
						return result, nil
					}
					if err := db.DB.Where("cluster_id = ?", clu.Cluster.ID).Preload("Host").Find(&nodes).Error; err != nil {
						ri.Result = StatusFailed
						ri.Msg = err.Error()
						result = append(result, ri)
						return result, nil
					}
					for _, node := range nodes {
						isExit := false
						for _, kn := range kubeNodes.Items {
							for _, addr := range kn.Status.Addresses {
								if addr.Type == "InternalIP" && node.Host.Ip == addr.Address {
									isExit = true
									continue
								}
							}
						}
						if !isExit {
							if err := db.DB.Model(&model.ClusterNode{}).Where("id = ?", node.ID).Updates(map[string]interface{}{"status": constant.StatusLost, "dirty": true}).Error; err != nil {
								ri.Result = StatusFailed
								ri.Msg = err.Error()
								result = append(result, ri)
								return result, nil
							}
						}
					}

					ri.Result = StatusSuccess
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
		msg := fmt.Sprintf("get cluster %s alive host falied: %s", c.Name, err.Error())
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

func getClusterToken(c model.Cluster) (string, error) {
	var master model.ClusterNode
	for _, item := range c.Nodes {
		if item.Role == constant.NodeRoleNameMaster {
			master = item
			break
		}
	}
	sshConfig := master.ToSSHConfig()
	client, _ := ssh.New(&sshConfig)
	return clusterUtil.GetClusterToken(client)
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
	client := &http.Client{Timeout: 2 * time.Second, Transport: tr}
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
