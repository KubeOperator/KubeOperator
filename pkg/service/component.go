package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
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
	clusterUtil "github.com/KubeOperator/KubeOperator/pkg/util/cluster"
	"github.com/jinzhu/gorm"
	"k8s.io/client-go/kubernetes"
)

const (
	npdPlaybook               = "12-npd.yml"
	metricServerPlaybook      = "13-metrics-server.yml"
	ingressControllerPlaybook = "14-ingress-controller.yml"
	gpuPlaybook               = "16-gpu-operator.yml"
	dnsCachePlaybook          = "17-dns-cache.yml"
	istioPlaybook             = "18-istio.yml"
)

type ComponentService interface {
	Get(clusterName string) ([]dto.Component, error)
	Create(component *dto.ComponentCreate) error
	Delete(clusterName, name string) error
	Sync(component *dto.ComponentSync) error
}

type componentService struct {
	clusterRepo    repository.ClusterRepository
	taskLogService TaskLogService
}

//  disable Initializing Waiting Failed enable Terminated

func NewComponentService() ComponentService {
	return &componentService{
		clusterRepo:    repository.NewClusterRepository(),
		taskLogService: NewTaskLogService(),
	}
}

func (c *componentService) Get(clusterName string) ([]dto.Component, error) {
	var (
		datas         []dto.Component
		dics          []model.ComponentDic
		specComponent []model.ClusterSpecComponent
		manifest      model.ClusterManifest
	)
	if err := db.DB.Find(&dics).Error; err != nil {
		return nil, err
	}
	cluster, err := c.clusterRepo.Get(clusterName)
	if err != nil {
		return nil, err
	}
	if err := db.DB.Where("name = ?", cluster.Version).Find(&manifest).Error; err != nil {
		return nil, err
	}
	var otherVars []dto.NameVersion
	if err := json.Unmarshal([]byte(manifest.OtherVars), &otherVars); err != nil {
		return nil, err
	}

	if err := db.DB.Where("cluster_id = ?", cluster.ID).Find(&specComponent).Error; err != nil {
		return nil, err
	}

	typeMap := make(map[string]bool)
	for _, dic := range dics {
		hasVersion := true
		for _, otherVar := range otherVars {
			if dic.Name == otherVar.Name {
				if dic.Version != otherVar.Version {
					hasVersion = false
				}
				break
			}
		}
		if !hasVersion {
			continue
		}
		data := dto.Component{
			Name:     dic.Name,
			Type:     dic.Type,
			Version:  dic.Version,
			Describe: dic.Describe,
		}
		isExit := false
		for _, spec := range specComponent {
			if disabled, ok := typeMap[spec.Type]; ok {
				if !disabled && spec.Status != constant.StatusDisabled {
					typeMap[spec.Type] = true
				}
			} else {
				if spec.Status != constant.StatusDisabled {
					typeMap[spec.Type] = true
				}
			}
			if dic.Name == spec.Name && dic.Version == spec.Version {
				isExit = true
				data.Status = spec.Status
				data.Message = spec.Message
				data.ID = spec.ID
				data.Vars = spec.Vars
				break
			}
		}
		if !isExit {
			data.Status = constant.StatusDisabled
			data.Message = ""
		}
		datas = append(datas, data)
	}
	for i := 0; i < len(datas); i++ {
		if disabled, ok := typeMap[datas[i].Type]; ok {
			datas[i].Disabled = disabled
		}
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
		if err := db.DB.Delete(&component).Error; err != nil {
			return err
		}
	}
	component = creation.ComponentCreate2Mo()
	component.ClusterID = cluster.ID
	component.Status = constant.StatusInitializing
	if err := db.DB.Create(&component).Error; err != nil {
		return err
	}

	playbook := c.loadPlayBookName(creation.Name)
	task := model.TaskLogDetail{
		ID:            component.ID,
		Name:          creation.Name,
		Task:          fmt.Sprintf("%s (%s)", playbook, constant.StatusEnabled),
		ClusterID:     cluster.ID,
		LastProbeTime: time.Now().Unix(),
		Status:        constant.TaskLogStatusRunning,
	}

	if err := c.taskLogService.StartDetail(&task); err != nil {
		c.errHandlerComponent(component, constant.StatusDisabled, err)
		return fmt.Errorf("save tasklog failed, err: %v", err)
	}

	writer, err := ansible.CreateAnsibleLogWriterWithId(cluster.Name, fmt.Sprintf("%s (%s)", component.ID, constant.StatusEnabled))
	if err != nil {
		_ = c.taskLogService.EndDetail(&task, component.Name, "component", constant.TaskLogStatusFailed, err.Error())
		return fmt.Errorf("create ansible log writer failed, err: %v", err)
	}

	if err := db.DB.Model(&model.ClusterSpecComponent{}).Where("id = ?", component.ID).
		Updates(map[string]interface{}{"status": constant.StatusInitializing, "message": ""}).Error; err != nil {
		_ = c.taskLogService.EndDetail(&task, component.Name, "component", constant.TaskLogStatusFailed, err.Error())
		return err
	}

	go c.docreate(&cluster, task, component, creation.Vars, writer)
	return nil
}

func (c componentService) docreate(cluster *model.Cluster, task model.TaskLogDetail, component model.ClusterSpecComponent, vars map[string]interface{}, writer io.Writer) {
	admCluster, err := c.loadAdmCluster(*cluster, component, vars, constant.StatusEnabled)
	if err != nil {
		_ = c.taskLogService.EndDetail(&task, component.Name, "component", constant.TaskLogStatusFailed, err.Error())
		c.errHandlerComponent(component, constant.StatusDisabled, err)
		return
	}

	client, err := clusterUtil.NewClusterClient(cluster)
	if err != nil {
		_ = c.taskLogService.EndDetail(&task, component.Name, "component", constant.TaskLogStatusFailed, err.Error())
		c.errHandlerComponent(component, constant.StatusDisabled, err)
		return
	}

	playbook := strings.ReplaceAll(task.Task, " (enable)", "")
	if err := phases.RunPlaybookAndGetResult(admCluster.Kobe, playbook, "", writer); err != nil {
		_ = c.taskLogService.EndDetail(&task, component.Name, "component", constant.TaskLogStatusFailed, err.Error())
		c.errHandlerComponent(component, constant.StatusFailed, err)
		return
	}
	_ = c.taskLogService.EndDetail(&task, component.Name, "component", constant.TaskLogStatusSuccess, "")
	component.Status = constant.StatusWaiting
	if err := db.DB.Save(&component).Error; err != nil {
		logger.Log.Errorf("save component status err: %s", err.Error())
		return
	}
	c.dosync([]model.ClusterSpecComponent{component}, client, component.Name)
}

func (c componentService) Delete(clusterName, name string) error {
	var component model.ClusterSpecComponent
	cluster, err := c.clusterRepo.GetWithPreload(clusterName, []string{"SpecConf", "SpecNetwork", "SpecRuntime", "Secret", "Nodes", "Nodes.Host", "Nodes.Host.Credential"})
	if err != nil {
		return err
	}
	db.DB.Where("name = ? AND cluster_id = ?", name, cluster.ID).First(&component)
	if component.ID == "" {
		return errors.New("not found")
	}

	playbook := c.loadPlayBookName(name)
	task := model.TaskLogDetail{
		ID:            fmt.Sprintf("%s (%s)", component.ID, constant.StatusDisabled),
		Name:          component.Name,
		Task:          fmt.Sprintf("%s (%s)", playbook, constant.StatusDisabled),
		ClusterID:     cluster.ID,
		LastProbeTime: time.Now().Unix(),
		Status:        constant.TaskLogStatusRunning,
	}
	if err := c.taskLogService.StartDetail(&task); err != nil {
		return fmt.Errorf("save tasklog failed, err: %v", err)
	}

	writer, err := ansible.CreateAnsibleLogWriterWithId(cluster.Name, fmt.Sprintf("%s (%s)", component.ID, constant.StatusDisabled))
	if err != nil {
		_ = c.taskLogService.EndDetail(&task, component.Name, "component", constant.TaskLogStatusFailed, err.Error())
		return fmt.Errorf("create ansible log writer failed, err: %v", err)
	}

	if err := db.DB.Model(&model.ClusterSpecComponent{}).Where("id = ?", component.ID).
		Updates(map[string]interface{}{"status": constant.StatusTerminating, "message": ""}).Error; err != nil {
		_ = c.taskLogService.EndDetail(&task, component.Name, "component", constant.TaskLogStatusFailed, err.Error())
		return err
	}

	go c.dodelete(&cluster, task, component, writer)

	return nil
}

func (c componentService) dodelete(cluster *model.Cluster, task model.TaskLogDetail, component model.ClusterSpecComponent, writer io.Writer) {
	admCluster, err := c.loadAdmCluster(*cluster, component, map[string]interface{}{}, constant.StatusDisabled)
	if err != nil {
		_ = c.taskLogService.EndDetail(&task, component.Name, "component", constant.TaskLogStatusFailed, err.Error())
		c.errHandlerComponent(component, constant.StatusFailed, err)
		return
	}
	playbook := strings.ReplaceAll(task.Task, " (disable)", "")
	if err := phases.RunPlaybookAndGetResult(admCluster.Kobe, playbook, "", writer); err != nil {
		_ = c.taskLogService.EndDetail(&task, component.Name, "component", constant.TaskLogStatusFailed, err.Error())
		c.errHandlerComponent(component, constant.StatusFailed, err)
		return
	}
	_ = c.taskLogService.EndDetail(&task, component.Name, "component", constant.TaskLogStatusSuccess, "")
	_ = db.DB.Where("id = ?", component.ID).Delete(&model.ClusterSpecComponent{})
}

func (c componentService) errHandlerComponent(component model.ClusterSpecComponent, status string, err error) {
	logger.Log.Errorf(err.Error())
	component.Status = status
	component.Message = err.Error()
	_ = db.DB.Save(&component)
}

func (c componentService) Sync(syncData *dto.ComponentSync) error {
	cluster, err := c.clusterRepo.GetWithPreload(syncData.ClusterName, []string{"SpecConf", "Secret", "Nodes", "Nodes.Host", "Nodes.Host.Credential"})
	if err != nil {
		return err
	}
	client, err := clusterUtil.NewClusterClient(&cluster)
	if err != nil {
		return err
	}
	var components []model.ClusterSpecComponent
	if err := db.DB.Where("cluster_id = ?", cluster.ID).Find(&components).Error; err != nil {
		return err
	}
	components, err = c.loadAllComponents(syncData.Names, cluster.ID, components)
	if err != nil {
		return err
	}

	if err := db.DB.Model(&model.ClusterSpecComponent{}).
		Where("cluster_id = ? AND name in (?)", cluster.ID, syncData.Names).
		Update("status", constant.StatusSynchronizing).Error; err != nil {
		return err
	}
	go c.dosync(components, client, syncData.Names...)

	return nil
}

func (c componentService) loadAllComponents(names []string, clusterID string, components []model.ClusterSpecComponent) ([]model.ClusterSpecComponent, error) {
	var dicList []model.ComponentDic
	if err := db.DB.Find(&dicList).Error; err != nil {
		return components, err
	}
	for _, name := range names {
		isExist := false
		for _, com := range components {
			if name == com.Name {
				isExist = true
				break
			}
		}
		if !isExist {
			for _, dic := range dicList {
				if dic.Name == name {
					comAdd := model.ClusterSpecComponent{
						Name:      dic.Name,
						ClusterID: clusterID,
						Version:   dic.Version,
						Type:      dic.Type,
						Status:    constant.StatusDisabled,
					}
					if err := db.DB.Create(&comAdd).Error; err != nil {
						fmt.Println(err)
					}
					components = append(components, comAdd)
				}
			}
		}
	}
	return components, nil
}

func (c componentService) dosync(components []model.ClusterSpecComponent, client *kubernetes.Clientset, names ...string) {
	for _, name := range names {
		switch name {
		case "gpu":
			if err := phases.WaitForDeployRunning("kube-operator", "gpu-operator", client); err != nil {
				c.changeStatus(components, name, constant.StatusFailed)
				continue
			}
			c.changeStatus(components, name, constant.StatusEnabled)
		case "ingress-nginx":
			if err := phases.WaitForDaemonsetRunning("kube-system", "nginx-ingress-controller", client); err != nil {
				if err := phases.WaitForDaemonsetRunning("kube-system", "ingress-nginx-controller", client); err != nil {
					c.changeStatus(components, name, constant.StatusFailed)
					continue
				}
			}
			c.changeStatus(components, name, constant.StatusEnabled)
		case "traefik":
			if err := phases.WaitForDaemonsetRunning("kube-system", "traefik", client); err != nil {
				c.changeStatus(components, name, constant.StatusFailed)
				continue
			}
			c.changeStatus(components, name, constant.StatusEnabled)
		case "dns-cache":
			if err := phases.WaitForDaemonsetRunning("kube-system", "node-local-dns", client); err != nil {
				c.changeStatus(components, name, constant.StatusFailed)
				continue
			}
			c.changeStatus(components, name, constant.StatusEnabled)
		case "metrics-server":
			if err := phases.WaitForDeployRunning("kube-system", "metrics-server", client); err != nil {
				c.changeStatus(components, name, constant.StatusFailed)
				continue
			}
			c.changeStatus(components, name, constant.StatusEnabled)
		case "npd":
			if err := phases.WaitForDaemonsetRunning("kube-system", "node-problem-detector", client); err != nil {
				c.changeStatus(components, name, constant.StatusFailed)
				continue
			}
			c.changeStatus(components, name, constant.StatusEnabled)
		case "istio":
			if err := phases.WaitForDeployRunning("istio-system", "istiod", client); err != nil {
				c.changeStatus(components, name, constant.StatusFailed)
				continue
			}
			c.changeStatus(components, name, constant.StatusEnabled)
		}
	}
}

func (c componentService) changeStatus(components []model.ClusterSpecComponent, name, status string) {
	for _, component := range components {
		if name == component.Name {
			if status == constant.StatusEnabled {
				component.Status = constant.StatusEnabled
				_ = db.DB.Save(component).Error
				continue
			}

			if status == constant.StatusFailed && component.Status == constant.StatusWaiting {
				component.Status = constant.StatusNotReady
				_ = db.DB.Save(component).Error
				continue
			}

			if status == constant.StatusFailed && component.Status == constant.StatusDisabled {
				_ = db.DB.Delete(component).Error
				continue
			}

			if status == constant.StatusFailed {
				component.Status = constant.StatusFailed
				component.Message = "can't found resource in cluster"
				_ = db.DB.Save(component).Error
				continue
			}
		}
	}
}

func (c componentService) loadAdmCluster(cluster model.Cluster, component model.ClusterSpecComponent, vars map[string]interface{}, operation string) (*adm.AnsibleHelper, error) {
	admCluster := adm.NewAnsibleHelper(cluster)

	if len(vars) != 0 {
		for k, v := range vars {
			if v != nil {
				admCluster.Kobe.SetVar(k, fmt.Sprintf("%v", v))
			}
		}
	}

	switch component.Name {
	case "gpu":
		admCluster.Kobe.SetVar(facts.SupportGpuFactName, operation)
	case "ingress-nginx":
		admCluster.Kobe.SetVar(facts.IngressControllerTypeFactName, "nginx")
		admCluster.Kobe.SetVar(facts.EnableNginxFactName, operation)
	case "traefik":
		admCluster.Kobe.SetVar(facts.IngressControllerTypeFactName, "traefik")
		admCluster.Kobe.SetVar(facts.EnableTraefikFactName, operation)
	case "dns-cache":
		admCluster.Kobe.SetVar(facts.EnableDnsCacheFactName, operation)
	case "metrics-server":
		admCluster.Kobe.SetVar(facts.MetricsServerFactName, operation)
	case "npd":
		admCluster.Kobe.SetVar(facts.EnableNpdFactName, operation)
	case "istio":
		admCluster.Kobe.SetVar(facts.EnableIstioFactName, operation)
	}
	return admCluster, nil
}

func (c componentService) loadPlayBookName(name string) string {
	switch name {
	case "gpu":
		return gpuPlaybook
	case "ingress-nginx", "traefik":
		return ingressControllerPlaybook
	case "dns-cache":
		return dnsCachePlaybook
	case "metrics-server":
		return metricServerPlaybook
	case "npd":
		return npdPlaybook
	case "istio":
		return istioPlaybook
	}
	return ""
}
