package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/logger"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/repository"
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm"
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/facts"
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/phases"
	"github.com/KubeOperator/KubeOperator/pkg/util/ansible"
	"github.com/jinzhu/gorm"
)

const (
	npdPlaybook               = "12-npd.yml"
	metricServerPlaybook      = "13-metrics-server.yml"
	ingressControllerPlaybook = "14-ingress-controller.yml"
	gpuPlaybook               = "16-gpu-operator.yml"
	dnsCachePlaybook          = "17-dns-cache.yml"
)

type ComponentService interface {
	Get(clusterName string) ([]dto.Component, error)
	Create(component *dto.ComponentCreate) error
}

type componentService struct {
	clusterRepo    repository.ClusterRepository
	taskLogService TaskLogService
	clusterService ClusterService
}

//  disable Initializing Waiting Failed enable Terminated

func NewComponentService() ComponentService {
	return &componentService{
		clusterRepo:    repository.NewClusterRepository(),
		taskLogService: NewTaskLogService(),
		clusterService: NewClusterService(),
	}
}

func (c *componentService) Get(clusterName string) ([]dto.Component, error) {
	var (
		datas         []dto.Component
		dics          []model.ComponentDic
		specComponent []model.ClusterSpecComponent
	)
	if err := db.DB.Find(&dics).Error; err != nil {
		return nil, err
	}
	cluster, err := c.clusterRepo.Get(clusterName)
	if err != nil {
		return nil, err
	}
	if err := db.DB.Where("cluster_id = ?", cluster.ID).Find(&specComponent).Error; err != nil {
		return nil, err
	}

	for _, dic := range dics {
		data := dto.Component{
			Name:     dic.Name,
			Version:  dic.Version,
			Describe: dic.Describe,
		}
		isExit := false
		for _, spec := range specComponent {
			if dic.Name == spec.Name && dic.Version == spec.Version {
				isExit = true
				data.Status = spec.Status
				data.Message = spec.Message
				data.ID = spec.ID
				break
			}
		}
		if !isExit {
			data.Status = constant.StatusDisabled
			data.Message = ""
		}
		datas = append(datas, data)
	}

	return datas, nil
}

func (c *componentService) Create(creation *dto.ComponentCreate) error {
	var component model.ClusterSpecComponent
	cluster, err := c.clusterRepo.GetWithPreload(creation.ClusterName, []string{"SpecConf", "SpecNetwork", "SpecRuntime", "Secret", "Nodes", "Nodes.Host", "Nodes.Host.Credential"})
	if err != nil {
		return err
	}
	if err := db.DB.Where("name = ? AND version = ? AND cluster_id = ?", creation.Name, creation.Version, cluster.ID).First(&component).Error; err != nil && !gorm.IsRecordNotFoundError(err) {
		return err
	}
	if component.ID != "" {
		if component.Status != constant.StatusDisabled && component.Status != constant.StatusFailed {
			return errors.New("COMPONENT_EXIST")
		}
	}
	if component.ID == "" {
		component = creation.ComponentCreate2Mo()
		component.ClusterID = cluster.ID
		if err := db.DB.Create(&component).Error; err != nil {
			return err
		}
	}
	playbook := c.loadPlayBookName(creation.Name)
	task := model.TaskLogDetail{
		ID:            component.ID,
		Task:          playbook,
		ClusterID:     cluster.ID,
		LastProbeTime: time.Now(),
		Status:        constant.TaskLogStatusRunning,
	}
	if err := c.taskLogService.StartDetail(&task); err != nil {
		return fmt.Errorf("save tasklog failed, err: %v", err)
	}

	//playbook
	go c.do(cluster, component, task)
	return nil
}

func (c componentService) do(cluster model.Cluster, component model.ClusterSpecComponent, task model.TaskLogDetail) {
	admCluster := adm.NewAnsibleHelper(cluster)
	writer, err := ansible.CreateAnsibleLogWriterWithId(cluster.Name, component.ID)
	if err != nil {
		logger.Log.Error(err)
	}
	if err := db.DB.Model(&model.ClusterSpecComponent{}).Where("id = ?", component.ID).Updates(map[string]interface{}{
		"status":  constant.StatusInitializing,
		"message": "",
	}).Error; err != nil {
		_ = c.taskLogService.EndDetail(&task, constant.TaskLogStatusFailed, err.Error())
		c.errCreateComponent(component, constant.StatusDisabled, err)
		return
	}

	// 获取 k8s client
	client, err := c.clusterService.NewClusterClient(cluster.Name)
	if err != nil {
		_ = c.taskLogService.EndDetail(&task, constant.TaskLogStatusFailed, err.Error())
		c.errCreateComponent(component, constant.StatusDisabled, err)
	}

	switch component.Name {
	case "gpu":
		admCluster.Kobe.SetVar(facts.SupportGpuFactName, constant.StatusEnabled)
		if err := phases.RunPlaybookAndGetResult(admCluster.Kobe, gpuPlaybook, "", writer); err != nil {
			_ = c.taskLogService.EndDetail(&task, constant.TaskLogStatusFailed, err.Error())
			c.errCreateComponent(component, constant.StatusFailed, fmt.Errorf("create component failed, err: %v", err.Error()))
			return
		}
		_ = c.taskLogService.EndDetail(&task, constant.TaskLogStatusSuccess, "")
		component.Status = constant.StatusWaiting
		if err := db.DB.Save(&component).Error; err != nil {
			logger.Log.Errorf("save component status err: %s", err.Error())
			return
		}
		if err := phases.WaitForDeployRunning("kube-operator", "gpu-operator", client); err != nil {
			c.errCreateComponent(component, constant.StatusNotReady, fmt.Errorf("waitting component running error %s", err.Error()))
			return
		}
	case "dns-cache":
		if err := phases.RunPlaybookAndGetResult(admCluster.Kobe, dnsCachePlaybook, "", writer); err != nil {
			_ = c.taskLogService.EndDetail(&task, constant.TaskLogStatusFailed, err.Error())
			c.errCreateComponent(component, constant.StatusFailed, fmt.Errorf("create component failed, err: %v", err.Error()))
			return
		}
		_ = c.taskLogService.EndDetail(&task, constant.TaskLogStatusSuccess, "")
		component.Status = constant.StatusWaiting
		if err := db.DB.Save(&component).Error; err != nil {
			logger.Log.Errorf("save component status err: %s", err.Error())
			return
		}
		if err := phases.WaitForDaemonsetRunning("kube-system", "node-local-dns", client); err != nil {
			c.errCreateComponent(component, constant.StatusNotReady, fmt.Errorf("waitting component running error %s", err.Error()))
			return
		}
	case "nginx":
		admCluster.Kobe.SetVar(facts.IngressControllerTypeFactName, "nginx")
		if err = phases.RunPlaybookAndGetResult(admCluster.Kobe, ingressControllerPlaybook, "", writer); err != nil {
			_ = c.taskLogService.EndDetail(&task, constant.TaskLogStatusFailed, err.Error())
			c.errCreateComponent(component, constant.StatusFailed, fmt.Errorf("create component failed, err: %v", err.Error()))
			return
		}
		_ = c.taskLogService.EndDetail(&task, constant.TaskLogStatusSuccess, "")
		component.Status = constant.StatusWaiting
		if err := db.DB.Save(&component).Error; err != nil {
			logger.Log.Errorf("save component status err: %s", err.Error())
			return
		}
		if err := phases.WaitForDaemonsetRunning("kube-system", "nginx-ingress-controller", client); err != nil {
			c.errCreateComponent(component, constant.StatusNotReady, fmt.Errorf("waitting component running error %s", err.Error()))
			return
		}
	case "traefik":
		admCluster.Kobe.SetVar(facts.IngressControllerTypeFactName, "traefik")
		if err = phases.RunPlaybookAndGetResult(admCluster.Kobe, ingressControllerPlaybook, "", writer); err != nil {
			_ = c.taskLogService.EndDetail(&task, constant.TaskLogStatusFailed, err.Error())
			c.errCreateComponent(component, constant.StatusFailed, fmt.Errorf("create component failed, err: %v", err.Error()))
			return
		}
		_ = c.taskLogService.EndDetail(&task, constant.TaskLogStatusSuccess, "")
		component.Status = constant.StatusWaiting
		if err := db.DB.Save(&component).Error; err != nil {
			logger.Log.Errorf("save component status err: %s", err.Error())
			return
		}
		if err := phases.WaitForDaemonsetRunning("kube-system", "traefik", client); err != nil {
			c.errCreateComponent(component, constant.StatusNotReady, fmt.Errorf("waitting component running error %s", err.Error()))
			return
		}
	case "metrics-server":
		admCluster.Kobe.SetVar(facts.MetricsServerFactName, constant.StatusEnabled)
		if err = phases.RunPlaybookAndGetResult(admCluster.Kobe, metricServerPlaybook, "", writer); err != nil {
			_ = c.taskLogService.EndDetail(&task, constant.TaskLogStatusFailed, err.Error())
			c.errCreateComponent(component, constant.StatusFailed, fmt.Errorf("create component failed, err: %v", err.Error()))
			return
		}
		_ = c.taskLogService.EndDetail(&task, constant.TaskLogStatusSuccess, "")
		component.Status = constant.StatusWaiting
		if err := db.DB.Save(&component).Error; err != nil {
			logger.Log.Errorf("save component status err: %s", err.Error())
			return
		}
		if err := phases.WaitForDeployRunning("kube-system", "metrics-server", client); err != nil {
			c.errCreateComponent(component, constant.StatusNotReady, fmt.Errorf("waitting component running error %s", err.Error()))
			return
		}
	case "npd":
		admCluster.Kobe.SetVar(facts.NpdFactName, constant.StatusEnabled)
		if err = phases.RunPlaybookAndGetResult(admCluster.Kobe, npdPlaybook, "", writer); err != nil {
			_ = c.taskLogService.EndDetail(&task, constant.TaskLogStatusFailed, err.Error())
			c.errCreateComponent(component, constant.StatusFailed, fmt.Errorf("create component failed, err: %v", err.Error()))
			return
		}
		_ = c.taskLogService.EndDetail(&task, constant.TaskLogStatusSuccess, "")
		component.Status = constant.StatusWaiting
		if err := db.DB.Save(&component).Error; err != nil {
			logger.Log.Errorf("save component status err: %s", err.Error())
			return
		}
		if err := phases.WaitForDaemonsetRunning("kube-system", "node-problem-detector", client); err != nil {
			c.errCreateComponent(component, constant.StatusNotReady, fmt.Errorf("waitting component running error %s", err.Error()))
			return
		}
	}
	component.Status = constant.StatusEnabled
	_ = db.DB.Save(&component)
}

func (c componentService) errCreateComponent(component model.ClusterSpecComponent, status string, err error) {
	logger.Log.Errorf(err.Error())
	component.Status = status
	component.Message = err.Error()
	_ = db.DB.Save(&component)
}

func (c componentService) loadPlayBookName(name string) string {
	switch name {
	case "gpu":
		return gpuPlaybook
	case "nginx", "traefik":
		return ingressControllerPlaybook
	case "dns-cache":
		return dnsCachePlaybook
	case "metrics-server":
		return metricServerPlaybook
	case "npd":
		return npdPlaybook
	}
	return ""
}
