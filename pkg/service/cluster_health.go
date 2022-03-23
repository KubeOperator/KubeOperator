package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/util/ipaddr"
	kubeUtil "github.com/KubeOperator/KubeOperator/pkg/util/kubernetes"
	"github.com/KubeOperator/KubeOperator/pkg/util/ssh"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
		checkKubernetesApiServer}
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
		go func(item int) {
			defer wg.Done()
			if err := ipaddr.Ping(c.Nodes[item].Host.Ip); err != nil {
				result.Level = constant.ClusterHealthLevelWarning
				result.Msg += err.Error()
				return
			}
			if c.Nodes[item].Role == constant.NodeRoleNameMaster {
				aliveMaster++
			}
		}(i)
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
		go func(item int) {
			defer wg.Done()
			sshCfg := c.Nodes[item].ToSSHConfig()
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
			if c.Nodes[item].Role == constant.NodeRoleNameMaster {
				aliveMaster++
			}
		}(i)
	}
	wg.Wait()
	if !(aliveMaster > 0) {
		result.Level = constant.ClusterHealthLevelError
	}
	return result
}
func checkKubernetesApiServer(c model.Cluster) dto.ClusterHealthHook {
	result := dto.ClusterHealthHook{
		Name:  HookNameCheckKubernetesApiConnection,
		Level: constant.ClusterHealthLevelSuccess,
	}
	clusterService := NewClusterService()
	secret, err := clusterService.GetSecrets(c.Name)
	if err != nil {
		result.Msg = fmt.Sprintf("get cluster %s alive host error %s", c.Name, err.Error())
		result.Level = constant.ClusterHealthLevelError
		return result
	}
	c.Secret = secret.ClusterSecret
	client, err := kubeUtil.NewKubernetesClient(&secret.KubeConf)
	if err != nil {
		result.Msg = fmt.Sprintf("get cluster %s client error %s", c.Name, err.Error())
		result.Level = constant.ClusterHealthLevelError
		return result
	}
	_, err = client.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		result.Msg = fmt.Sprintf("get cluster %s api error %s", c.Name, err.Error())
		result.Level = constant.ClusterHealthLevelError
		return result
	}
	return result
}

var resolveMethods = map[string]string{
	HookNameCheckHostsNetwork:            "No method",
	HookNameCheckHostsSSHConnection:      "No method",
	HookNameCheckKubernetesApiConnection: "Get kubernetes token again",
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
				}
			}
		}
	default:
		return result, nil
	}
	return result, nil
}
