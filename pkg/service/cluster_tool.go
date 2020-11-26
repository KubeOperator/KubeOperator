package service

import (
	"context"
	"encoding/json"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/repository"
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/tools"
	kubernetesUtil "github.com/KubeOperator/KubeOperator/pkg/util/kubernetes"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ClusterToolService interface {
	List(clusterName string) ([]dto.ClusterTool, error)
	Enable(clusterName string, tool dto.ClusterTool) (dto.ClusterTool, error)
	Disable(clusterName string, tool dto.ClusterTool) (dto.ClusterTool, error)
}

func NewClusterToolService() ClusterToolService {
	return &clusterToolService{
		toolRepo:       repository.NewClusterToolRepository(),
		clusterService: NewClusterService(),
	}
}

type clusterToolService struct {
	toolRepo       repository.ClusterToolRepository
	clusterService ClusterService
}

func (c clusterToolService) List(clusterName string) ([]dto.ClusterTool, error) {
	var items []dto.ClusterTool
	ms, err := c.toolRepo.List(clusterName)
	if err != nil {
		return items, err
	}
	for _, m := range ms {
		d := dto.ClusterTool{ClusterTool: m}
		d.Vars = map[string]interface{}{}
		_ = json.Unmarshal([]byte(m.Vars), &d.Vars)
		items = append(items, d)
	}
	return items, nil
}

func (c clusterToolService) Disable(clusterName string, tool dto.ClusterTool) (dto.ClusterTool, error) {
	cluster, err := c.clusterService.Get(clusterName)
	if err != nil {
		return tool, err
	}
	tool.ClusterID = cluster.ID
	mo := tool.ClusterTool
	buf, _ := json.Marshal(&tool.Vars)
	mo.Vars = string(buf)
	tool.ClusterTool = mo
	endpoint, err := c.clusterService.GetApiServerEndpoint(clusterName)
	if err != nil {
		return tool, err
	}
	secret, err := c.clusterService.GetSecrets(clusterName)
	if err != nil {
		return tool, err
	}

	itemValue, ok := tool.Vars["namespace"]
	namespace := ""
	if !ok {
		namespace = constant.DefaultNamespace
	} else {
		namespace = itemValue.(string)
	}

	ct, err := tools.NewClusterTool(&tool.ClusterTool, cluster.Cluster, endpoint, secret.ClusterSecret, namespace)
	if err != nil {
		return tool, err
	}
	mo.Status = constant.ClusterTerminating
	_ = c.toolRepo.Save(&mo)
	go c.doUninstall(ct, &tool.ClusterTool)
	return tool, nil
}

func (c clusterToolService) Enable(clusterName string, tool dto.ClusterTool) (dto.ClusterTool, error) {
	cluster, err := c.clusterService.Get(clusterName)
	if err != nil {
		return tool, err
	}
	tool.ClusterID = cluster.ID
	mo := tool.ClusterTool
	buf, _ := json.Marshal(&tool.Vars)
	mo.Vars = string(buf)
	tool.ClusterTool = mo
	endpoint, err := c.clusterService.GetApiServerEndpoint(clusterName)
	if err != nil {
		return tool, err
	}
	secret, err := c.clusterService.GetSecrets(clusterName)
	if err != nil {
		return tool, err
	}

	kubeClient, err := kubernetesUtil.NewKubernetesClient(&kubernetesUtil.Config{
		Host:  endpoint.Address,
		Token: secret.KubernetesToken,
		Port:  endpoint.Port,
	})
	if err != nil {
		return tool, err
	}

	itemValue, ok := tool.Vars["namespace"]
	namespace := ""
	if !ok {
		namespace = constant.DefaultNamespace
	} else {
		namespace = itemValue.(string)
	}

	ns, _ := kubeClient.CoreV1().Namespaces().Get(context.TODO(), namespace, metav1.GetOptions{})
	if ns.ObjectMeta.Name == "" {
		n := &v1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: namespace,
			},
		}
		_, err = kubeClient.CoreV1().Namespaces().Create(context.TODO(), n, metav1.CreateOptions{})
		if err != nil {
			return tool, err
		}
	}
	ct, err := tools.NewClusterTool(&tool.ClusterTool, cluster.Cluster, endpoint, secret.ClusterSecret, namespace)
	if err != nil {
		return tool, err
	}
	mo.Status = constant.ClusterInitializing
	_ = c.toolRepo.Save(&mo)
	go c.doInstall(ct, &tool.ClusterTool)
	return tool, nil
}

func (c clusterToolService) doInstall(p tools.Interface, tool *model.ClusterTool) {
	err := p.Install()
	if err != nil {
		tool.Status = constant.ClusterFailed
		tool.Message = err.Error()
	} else {
		tool.Status = constant.ClusterRunning
	}
	_ = c.toolRepo.Save(tool)
}

func (c clusterToolService) doUninstall(p tools.Interface, tool *model.ClusterTool) {
	_ = p.Uninstall()
	tool.Status = constant.ClusterWaiting
	_ = c.toolRepo.Save(tool)
}
