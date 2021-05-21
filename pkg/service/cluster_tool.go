package service

import (
	"context"
	"encoding/json"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/logger"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/repository"
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/tools"
	"github.com/KubeOperator/KubeOperator/pkg/util/kubernetes"
	kubernetesUtil "github.com/KubeOperator/KubeOperator/pkg/util/kubernetes"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ClusterToolService interface {
	List(clusterName string) ([]dto.ClusterTool, error)
	Enable(clusterName string, tool dto.ClusterTool) (dto.ClusterTool, error)
	Upgrade(clusterName string, tool dto.ClusterTool) (dto.ClusterTool, error)
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
	cluster, hosts, secret, err := c.getBaseParams(clusterName)
	if err != nil {
		return tool, err
	}

	tool.ClusterID = cluster.ID
	mo := tool.ClusterTool
	buf, _ := json.Marshal(&tool.Vars)
	mo.Vars = string(buf)
	tool.ClusterTool = mo

	itemValue, ok := tool.Vars["namespace"]
	namespace := ""
	if !ok {
		namespace = constant.DefaultNamespace
	} else {
		namespace = itemValue.(string)
	}

	ct, err := tools.NewClusterTool(&tool.ClusterTool, cluster.Cluster, hosts, secret.ClusterSecret, namespace, namespace, false)
	if err != nil {
		return tool, err
	}
	mo.Status = constant.ClusterTerminating
	_ = c.toolRepo.Save(&mo)
	go c.doUninstall(ct, &tool.ClusterTool)
	return tool, nil
}

func (c clusterToolService) Enable(clusterName string, tool dto.ClusterTool) (dto.ClusterTool, error) {
	cluster, hosts, secret, err := c.getBaseParams(clusterName)
	if err != nil {
		return tool, err
	}

	var toolDetail model.ClusterToolDetail
	if err := db.DB.Where("name = ? AND version = ?", tool.Name, tool.Version).Find(&toolDetail).Error; err != nil {
		return tool, err
	}

	tool.ClusterID = cluster.ID
	mo := tool.ClusterTool
	buf, _ := json.Marshal(&tool.Vars)
	mo.Vars = string(buf)
	tool.ClusterTool = mo

	kubeClient, err := kubernetesUtil.NewKubernetesClient(&kubernetesUtil.Config{
		Hosts: hosts,
		Token: secret.KubernetesToken,
	})
	if err != nil {
		return tool, err
	}
	oldNamespace, namespace := c.getNamespace(cluster.ID, tool)
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
	ct, err := tools.NewClusterTool(&tool.ClusterTool, cluster.Cluster, hosts, secret.ClusterSecret, oldNamespace, namespace, true)
	if err != nil {
		return tool, err
	}
	mo.Status = constant.ClusterInitializing
	_ = c.toolRepo.Save(&mo)
	go c.doInstall(ct, &tool.ClusterTool, toolDetail)
	return tool, nil
}

func (c clusterToolService) Upgrade(clusterName string, tool dto.ClusterTool) (dto.ClusterTool, error) {
	cluster, hosts, secret, err := c.getBaseParams(clusterName)
	if err != nil {
		return tool, err
	}

	var toolDetail model.ClusterToolDetail
	if err := db.DB.Where("name = ? AND version = ?", tool.Name, tool.HigherVersion).Find(&toolDetail).Error; err != nil {
		return tool, err
	}

	tool.ClusterID = cluster.ID
	mo := tool.ClusterTool
	buf, _ := json.Marshal(&tool.Vars)
	mo.Vars = string(buf)
	mo.Status = constant.ClusterUpgrading
	mo.Version = mo.HigherVersion
	mo.HigherVersion = ""
	tool.ClusterTool = mo

	itemValue, ok := tool.Vars["namespace"]
	namespace := ""
	if !ok {
		namespace = constant.DefaultNamespace
	} else {
		namespace = itemValue.(string)
	}
	ct, err := tools.NewClusterTool(&tool.ClusterTool, cluster.Cluster, hosts, secret.ClusterSecret, namespace, namespace, true)
	if err != nil {
		return tool, err
	}

	_ = c.toolRepo.Save(&mo)
	go c.doUpgrade(ct, &tool.ClusterTool, toolDetail)
	return tool, nil
}

func (c clusterToolService) doInstall(p tools.Interface, tool *model.ClusterTool, toolDetail model.ClusterToolDetail) {
	err := p.Install(toolDetail)
	if err != nil {
		tool.Status = constant.ClusterFailed
		tool.Message = err.Error()
	} else {
		tool.Status = constant.ClusterRunning
	}
	_ = c.toolRepo.Save(tool)
}

func (c clusterToolService) doUpgrade(p tools.Interface, tool *model.ClusterTool, toolDetail model.ClusterToolDetail) {
	err := p.Upgrade(toolDetail)
	if err != nil {
		tool.Status = constant.ClusterFailed
		tool.Message = err.Error()
	} else {
		tool.Status = constant.ClusterRunning
	}
	_ = c.toolRepo.Save(tool)
}

func (c clusterToolService) doUninstall(p tools.Interface, tool *model.ClusterTool) {
	if err := p.Uninstall(); err != nil {
		logger.Log.Errorf("do uninstall tool-%s failed, error: %s", tool.Name, err.Error())
	}
	tool.Status = constant.ClusterWaiting
	_ = c.toolRepo.Save(tool)
}

func (c clusterToolService) getNamespace(clusterID string, tool dto.ClusterTool) (string, string) {
	namespace := ""
	Sp, ok := tool.Vars["namespace"]
	if !ok {
		namespace = constant.DefaultNamespace
	} else {
		namespace = Sp.(string)
	}
	var oldTools model.ClusterTool
	if err := db.DB.Where("cluster_id = ? AND name = ?", clusterID, tool.Name).First(&oldTools).Error; err != nil {
		return namespace, namespace
	}
	oldVars := map[string]interface{}{}
	_ = json.Unmarshal([]byte(oldTools.Vars), &oldVars)
	oldSp, ok := oldVars["namespace"]
	if !ok {
		return namespace, namespace
	} else {
		return oldSp.(string), namespace
	}
}

func (c clusterToolService) getBaseParams(clusterName string) (dto.Cluster, []kubernetes.Host, dto.ClusterSecret, error) {
	var (
		cluster dto.Cluster
		host    []kubernetes.Host
		secret  dto.ClusterSecret
		err     error
	)
	if err := db.DB.Where("name = ?", clusterName).Preload("Spec").Find(&cluster).Error; err != nil {
		return cluster, host, secret, err
	}

	host, err = c.clusterService.GetApiServerEndpoints(clusterName)
	if err != nil {
		return cluster, host, secret, err
	}
	secret, err = c.clusterService.GetSecrets(clusterName)
	if err != nil {
		return cluster, host, secret, err
	}

	return cluster, host, secret, nil
}
