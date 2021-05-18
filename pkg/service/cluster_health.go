package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/util/ipaddr"
	kubeUtil "github.com/KubeOperator/KubeOperator/pkg/util/kubernetes"
	"github.com/KubeOperator/KubeOperator/pkg/util/ssh"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
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
		checkHostNetworkConnected,
		checkHostSSHConnected,
		checkKubernetesApiServer,
		checkKubernetesNodeStatus,
	}
	clu, err := c.clusterService.Get(clusterName)
	if err != nil {
		return nil, err
	}
	result := dto.ClusterHealth{}
	for _, h := range hookList {
		hookResult := h(clu.Cluster)
		if hookResult.Level == constant.ClusterHealthLevelError {
			result.Level = constant.ClusterHealthLevelError
		}
		result.Hooks = append(result.Hooks, hookResult)
	}
	return &result, nil
}

const (
	HookNameCheckHostsNetwork            = "Check Hosts Network"
	HookNameCheckHostsSSHConnection      = "Check Host SSH Connection"
	HookNameCheckKubernetesApiConnection = "Check Kubernetes Api Connection"
	HookNameCheckKubernetesNodeStatus    = "Check Kubernetes Node Status"
)

func checkHostNetworkConnected(c model.Cluster) dto.ClusterHealthHook {
	result := dto.ClusterHealthHook{
		Name:  HookNameCheckHostsNetwork,
		Level: constant.ClusterHealthLevelSuccess,
	}
	aliveMaster := 0
	wg := &sync.WaitGroup{}
	for i := range c.Nodes {
		wg.Add(1)
		i := i
		go func() {
			defer wg.Done()
			if err := ipaddr.Ping(c.Nodes[i].Host.Ip); err != nil {
				result.Level = constant.ClusterHealthLevelWarning
				result.Msg += err.Error()
				return
			}
			if c.Nodes[i].Role == constant.NodeRoleNameMaster {
				aliveMaster++
			}
		}()
	}
	wg.Wait()
	if !(aliveMaster > 0) {
		result.Level = constant.ClusterHealthLevelError
	}
	return result
}
func checkHostSSHConnected(c model.Cluster) dto.ClusterHealthHook {
	result := dto.ClusterHealthHook{
		Name:  HookNameCheckHostsSSHConnection,
		Level: constant.ClusterHealthLevelSuccess,
	}
	aliveMaster := 0
	wg := &sync.WaitGroup{}
	for i := range c.Nodes {
		wg.Add(1)
		i := i
		go func() {
			defer wg.Done()
			sshCfg := c.Nodes[i].ToSSHConfig()
			sshClient, err := ssh.New(&sshCfg)
			if err != nil {
				result.Msg += err.Error()
				return
			}
			if err := sshClient.Ping(); err != nil {
				result.Level = constant.ClusterHealthLevelWarning
				result.Msg += err.Error()
				return
			}
			if c.Nodes[i].Role == constant.NodeRoleNameMaster {
				aliveMaster++
			}
		}()
	}
	wg.Wait()
	if !(aliveMaster > 0) {
		result.Level = constant.ClusterHealthLevelError
	}
	return result
}
func checkKubernetesApiServer(c model.Cluster) dto.ClusterHealthHook {
	client, level, msg := getBaseParams(c, HookNameCheckKubernetesNodeStatus)
	result := dto.ClusterHealthHook{
		Name:  HookNameCheckKubernetesNodeStatus,
		Level: level,
		Msg:   msg,
	}
	if len(msg) != 0 {
		return result
	}

	_, err := client.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		result.Msg = fmt.Sprintf("get cluster %s api error %s", c.Name, err.Error())
		result.Level = constant.ClusterHealthLevelError
		return result
	}
	return result
}

func checkKubernetesNodeStatus(c model.Cluster) dto.ClusterHealthHook {
	var nodes []model.ClusterNode
	client, level, msg := getBaseParams(c, HookNameCheckKubernetesNodeStatus)
	result := dto.ClusterHealthHook{
		Name:  HookNameCheckKubernetesNodeStatus,
		Level: level,
		Msg:   msg,
	}
	if len(msg) != 0 {
		return result
	}

	kubeNodes, err := client.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		result.Msg = fmt.Sprintf("get cluster %s kubeNodes error %s", c.Name, err.Error())
		result.Level = constant.ClusterHealthLevelError
		return result
	}
	if err := db.DB.Where("cluster_id = ?", c.ID).Find(&nodes).Error; err != nil {
		result.Msg = fmt.Sprintf("get cluster %s nodes from db error %s", c.Name, err.Error())
		result.Level = constant.ClusterHealthLevelError
		return result
	}
	if len(nodes) != len(kubeNodes.Items) {
		result.Msg = fmt.Sprintf("The number of system nodes does not match the number of k8s nodes")
		result.Level = constant.ClusterHealthLevelError
		return result
	}

	return result
}

var resolveMethods = map[string]string{
	HookNameCheckHostsNetwork:            "No method",
	HookNameCheckHostsSSHConnection:      "No method",
	HookNameCheckKubernetesApiConnection: "Get kubernetes token again",
	HookNameCheckKubernetesNodeStatus:    "Update cluster node status",
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
	case constant.ClusterHealthLevelError:
		for i := range ch.Hooks {
			if ch.Hooks[i].Level == constant.ClusterHealthLevelError {
				switch ch.Hooks[i].Name {
				case HookNameCheckHostsNetwork, HookNameCheckHostsSSHConnection:
					ri := dto.ClusterRecoverItem{
						Name:     resolveMethods[ch.Hooks[i].Name],
						HookName: ch.Hooks[i].Name,
						Result:   constant.StatusFailed,
						Msg:      "No method",
					}
					result = append(result, ri)
					return result, nil
				case HookNameCheckKubernetesApiConnection:
					ri := dto.ClusterRecoverItem{
						Name:     resolveMethods[ch.Hooks[i].Name],
						HookName: ch.Hooks[i].Name,
					}
					err := c.clusterInitService.GatherKubernetesToken(clu.Cluster)
					if err != nil {
						ri.Result = constant.StatusFailed
						ri.Msg = err.Error()
						result = append(result, ri)
						return result, nil
					}
					ri.Result = constant.StatusSuccess
					result = append(result, ri)
				case HookNameCheckKubernetesNodeStatus:
					client, _, msg := getBaseParams(clu.Cluster, HookNameCheckKubernetesNodeStatus)
					ri := dto.ClusterRecoverItem{
						Name:     resolveMethods[ch.Hooks[i].Name],
						HookName: ch.Hooks[i].Name,
					}
					if len(msg) != 0 {
						ri.Result = constant.StatusFailed
						ri.Msg = err.Error()
						result = append(result, ri)
						return result, nil
					}

					var nodes []model.ClusterNode
					kubeNodes, err := client.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
					if err != nil {
						ri.Result = constant.StatusFailed
						ri.Msg = err.Error()
						result = append(result, ri)
						return result, nil
					}
					if err := db.DB.Where("cluster_id = ?", clu.Cluster.ID).Preload("Host").Find(&nodes).Error; err != nil {
						ri.Result = constant.StatusFailed
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
								ri.Result = constant.StatusFailed
								ri.Msg = err.Error()
								result = append(result, ri)
								return result, nil
							}
						}
					}

					ri.Result = constant.StatusSuccess
					result = append(result, ri)
				default:
					return result, nil
				}
			}
		}
	}

	return result, nil
}

func getBaseParams(c model.Cluster, name string) (*kubernetes.Clientset, string, string) {
	var clusterService = NewClusterService()
	secret, err := clusterService.GetSecrets(c.Name)
	if err != nil {
		msg := fmt.Sprintf("get cluster %s secret error %s", c.Name, err.Error())
		level := constant.ClusterHealthLevelError
		return nil, level, msg
	}

	endpoints, err := clusterService.GetApiServerEndpoints(c.Name)
	if err != nil {
		msg := fmt.Sprintf("get cluster %s endpoint error %s", c.Name, err.Error())
		level := constant.ClusterHealthLevelError
		return nil, level, msg
	}

	_, err = kubeUtil.SelectAliveHost(endpoints)
	if err != nil {
		msg := fmt.Sprintf("get cluster %s alive host error %s", c.Name, err.Error())
		level := constant.ClusterHealthLevelError
		return nil, level, msg
	}

	kubeClient, err := kubeUtil.NewKubernetesClient(&kubeUtil.Config{
		Hosts: endpoints,
		Token: secret.KubernetesToken,
	})
	if err != nil {
		msg := fmt.Sprintf("get cluster %s kubeclient error %s", c.Name, err.Error())
		level := constant.ClusterHealthLevelError
		return nil, level, msg
	}

	return kubeClient, constant.ClusterHealthLevelSuccess, ""
}
