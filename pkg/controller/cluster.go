package controller

import (
	"encoding/json"
	"errors"
	"fmt"

	"gopkg.in/yaml.v3"

	"github.com/KubeOperator/KubeOperator/pkg/controller/condition"
	"github.com/KubeOperator/KubeOperator/pkg/logger"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/controller/kolog"
	"github.com/KubeOperator/KubeOperator/pkg/controller/page"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/service"
	"github.com/kataras/iris/v12/context"
)

type ClusterController struct {
	Ctx                   context.Context
	ClusterService        service.ClusterService
	ClusterInitService    service.ClusterInitService
	ClusterNodeService    service.ClusterNodeService
	ClusterImportService  service.ClusterImportService
	CisService            service.CisService
	ClusterUpgradeService service.ClusterUpgradeService
	ClusterHealthService  service.ClusterHealthService
	BackupAccountService  service.BackupAccountService
}

func NewClusterController() *ClusterController {
	return &ClusterController{
		ClusterService:        service.NewClusterService(),
		ClusterInitService:    service.NewClusterInitService(),
		ClusterNodeService:    service.NewClusterNodeService(),
		ClusterImportService:  service.NewClusterImportService(),
		CisService:            service.NewCisService(),
		ClusterUpgradeService: service.NewClusterUpgradeService(),
		ClusterHealthService:  service.NewClusterHealthService(),
		BackupAccountService:  service.NewBackupAccountService(),
	}
}

// List Cluster
// @Tags clusters
// @Summary Show all clusters
// @Description Show clusters
// @Accept  json
// @Produce  json
// @Success 200 {object} page.Page
// @Security ApiKeyAuth
// @Router /clusters/ [get]
func (c ClusterController) Get() (*dto.ClusterPage, error) {
	page, _ := c.Ctx.Values().GetBool("page")
	var conditions condition.Conditions
	sessionUser := c.Ctx.Values().Get("user")
	user, _ := sessionUser.(dto.SessionUser)
	if c.Ctx.GetContentLength() > 0 {
		if err := c.Ctx.ReadJSON(&conditions); err != nil {
			return nil, err
		}
	}
	if page {
		num, _ := c.Ctx.Values().GetInt(constant.PageNumQueryKey)
		size, _ := c.Ctx.Values().GetInt(constant.PageSizeQueryKey)
		isPolling := c.Ctx.URLParam("isPolling")
		pageItem, err := c.ClusterService.Page(num, size, isPolling, user, conditions)
		if err != nil {
			return nil, err
		}
		return pageItem, nil
	} else {
		var pageItem dto.ClusterPage
		items, err := c.ClusterService.List()
		if err != nil {
			return nil, err
		}
		pageItem.Items = items
		pageItem.Total = len(items)
		return &pageItem, nil
	}
}

// Search Cluster
// @Tags clusters
// @Summary Search cluster
// @Description 过滤集群
// @Accept  json
// @Produce  json
// @Param conditions body condition.Conditions true "conditions"
// @Success 200 {object} page.Page
// @Security ApiKeyAuth
// @Router /clusters/search [post]
func (c ClusterController) PostSearch() (*dto.ClusterPage, error) {
	page, _ := c.Ctx.Values().GetBool("page")
	var conditions condition.Conditions
	sessionUser := c.Ctx.Values().Get("user")
	user, _ := sessionUser.(dto.SessionUser)
	if c.Ctx.GetContentLength() > 0 {
		if err := c.Ctx.ReadJSON(&conditions); err != nil {
			return nil, err
		}
	}
	if page {
		num, _ := c.Ctx.Values().GetInt(constant.PageNumQueryKey)
		size, _ := c.Ctx.Values().GetInt(constant.PageSizeQueryKey)
		isPolling := c.Ctx.URLParam("isPolling")
		pageItem, err := c.ClusterService.Page(num, size, isPolling, user, conditions)
		if err != nil {
			return nil, err
		}
		return pageItem, nil
	} else {
		var pageItem dto.ClusterPage
		items, err := c.ClusterService.List()
		if err != nil {
			return nil, err
		}
		pageItem.Items = items
		pageItem.Total = len(items)
		return &pageItem, nil
	}
}

// Get Cluster
// @Tags clusters
// @Summary Show a cluster
// @Description Show a cluster
// @Accept  json
// @Produce  json
// @Success 200 {object} dto.Cluster
// @Security ApiKeyAuth
// @Router /clusters/{name}/ [get]
func (c ClusterController) GetBy(name string) (*dto.Cluster, error) {
	cl, err := c.ClusterService.Get(name)
	if err != nil {
		return nil, err
	}
	return &cl, nil
}

// Get Cluster Name By Projects
// @Tags clusters
// @Summary Show cluster names of projects
// @Description Show a cluster names of projects
// @Accept  json
// @Produce  json
// @Success 200 {object} dto.Cluster
// @Security ApiKeyAuth
// @Router /clusters/name/{projectNames} [get]
func (c ClusterController) GetNameBy(projectNames string) ([]dto.ClusterInfo, error) {
	return c.ClusterService.GetClusterByProject(projectNames)
}

func (c ClusterController) GetExistenceBy(name string) *dto.IsClusterNameExist {
	isExit := c.ClusterService.CheckExistence(name)
	return &dto.IsClusterNameExist{IsExist: isExit}
}

func (c ClusterController) GetStatusBy(name string) (*dto.TaskLog, error) {
	return c.ClusterService.GetStatus(name)
}

// Create Cluster
// @Tags clusters
// @Summary Create a cluster
// @Description Create a cluster
// @Param request body dto.ClusterCreate true "request"
// @Accept  json
// @Produce  json
// @Success 200 {object} dto.Cluster
// @Security ApiKeyAuth
// @Router /clusters/ [post]
func (c ClusterController) Post() (*dto.Cluster, error) {
	var req dto.ClusterCreate
	err := c.Ctx.ReadJSON(&req)
	if err != nil {
		return nil, err
	}
	item, err := c.ClusterService.Create(req)
	if err != nil {
		logger.Log.Info(fmt.Sprintf("%+v", err))
		return nil, err
	}
	operator := c.Ctx.Values().GetString("operator")
	go kolog.Save(operator, constant.CREATE_CLUSTER, req.Name)
	return item, nil
}

func (c ClusterController) PostInitBy(name string) error {
	operator := c.Ctx.Values().GetString("operator")
	go kolog.Save(operator, constant.INIT_CLUSTER, name)

	return c.ClusterService.ReCreate(name)
}

// Load Cluster Info for import
// @Tags clusters
// @Summary Load cluster info
// @Description Upgrade a cluster
// @Param request body dto.ClusterLoad true "request"
// @Accept  json
// @Produce  json
// @Success 200 {object} dto.ClusterLoadInfo
// @Security ApiKeyAuth
// @Router /clusters/load [post]
func (c ClusterController) PostLoad() (dto.ClusterLoadInfo, error) {
	var req dto.ClusterLoad
	var data dto.ClusterLoadInfo
	err := c.Ctx.ReadJSON(&req)
	if err != nil {
		return data, err
	}

	data, err = c.ClusterImportService.LoadClusterInfo(&req)
	if err != nil {
		return data, err
	}
	return data, nil
}

// Upgrade Cluster
// @Tags clusters
// @Summary Upgrade a cluster
// @Description Upgrade a cluster
// @Param request body dto.ClusterUpgrade true "request"
// @Accept  json
// @Produce  json
// @Success 200
// @Security ApiKeyAuth
// @Router /clusters/upgrade [post]
func (c ClusterController) PostUpgrade() error {
	var req dto.ClusterUpgrade
	err := c.Ctx.ReadJSON(&req)
	if err != nil {
		logger.Log.Info(fmt.Sprintf("%+v", err))
		return err
	}

	operator := c.Ctx.Values().GetString("operator")
	go kolog.Save(operator, constant.UPGRADE_CLUSTER, req.ClusterName+"("+req.Version+")")

	return c.ClusterUpgradeService.Upgrade(req)
}

// Delete Cluster
// @Tags clusters
// @Summary Delete a cluster
// @Description delete a cluster by name
// @Param force query string true  "是否强制（true, false）"
// @Param uninstall query string true  "是否卸载（true, false）"
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Router /clusters/{name}/ [delete]
func (c ClusterController) DeleteBy(name string) error {
	operator := c.Ctx.Values().GetString("operator")
	force, _ := c.Ctx.Values().GetBool("force")
	uninstallStr := c.Ctx.URLParam("uninstall")
	uninstall := uninstallStr == "true"

	go kolog.Save(operator, constant.DELETE_CLUSTER, name)
	return c.ClusterService.Delete(name, force, uninstall)
}

// Import Cluster
// @Tags clusters
// @Summary Import a cluster
// @Description import a cluster
// @Param request body dto.ClusterImport true "request"
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Router /clusters/import/ [post]
func (c ClusterController) PostImport() error {
	var req dto.ClusterImport
	err := c.Ctx.ReadJSON(&req)
	if err != nil {
		return err
	}

	operator := c.Ctx.Values().GetString("operator")
	go kolog.Save(operator, constant.IMPORT_CLUSTER, req.Name)

	return c.ClusterImportService.Import(req)
}

// Get Cluster Nodes
// @Tags clusters
// @Summary Get cluster nodes
// @Description Get cluster nodes
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Router /clusters/node/{clusterName} [get]
func (c ClusterController) GetNodeBy(clusterName string) (*dto.NodePage, error) {
	p, _ := c.Ctx.Values().GetBool("page")
	if p {
		num, _ := c.Ctx.Values().GetInt(constant.PageNumQueryKey)
		size, _ := c.Ctx.Values().GetInt(constant.PageSizeQueryKey)
		isPolling := c.Ctx.URLParam("isPolling")
		pageItem, err := c.ClusterNodeService.Page(num, size, isPolling, clusterName)
		if err != nil {
			return nil, err
		}
		return pageItem, nil
	} else {
		var pageItem dto.NodePage
		cns, err := c.ClusterNodeService.List(clusterName)
		if err != nil {
			return nil, err
		}
		pageItem.Items = cns
		pageItem.Total = len(cns)
		return &pageItem, nil
	}

}

// Get Cluster Details
// @Tags clusters
// @Summary Get cluster node details
// @Description Get cluster node details
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {object} dto.Node
// @Router /clusters/node/{clusterName}/{nodeName} [get]
func (c ClusterController) GetNodeDetailBy(clusterName string, nodeName string) (*dto.Node, error) {
	node, err := c.ClusterNodeService.Get(clusterName, nodeName)
	if err != nil {
		return nil, err
	}
	return node, nil
}

// Recreate Cluster Node
// @Tags clusters
// @Summary Recreate cluster node
// @Description Recreate cluster node
// @Accept  json
// @Produce  json
// @Param request body dto.NodeBatch true "request"
// @Security ApiKeyAuth
// @Success 200
// @Router /clusters/node/recreate/{clusterName} [post]
func (c ClusterController) PostNodeRecreateBy(clusterName string) error {
	var req dto.NodeBatch
	err := c.Ctx.ReadJSON(&req)
	if err != nil {
		return err
	}
	if err := c.ClusterNodeService.Recreate(clusterName, req); err != nil {
		return err
	}
	return nil
}

// Batch Delete Or Create Cluster Node
// @Tags clusters
// @Summary Batch delete or create cluster node
// @Description Batch delete or create cluster node
// @Accept  json
// @Produce  json
// @Param request body dto.NodeBatch true "request"
// @Success 200
// @Security ApiKeyAuth
// @Router /clusters/node/batch/{clusterName} [post]
func (c ClusterController) PostNodeBatchBy(clusterName string) error {
	var req dto.NodeBatch
	err := c.Ctx.ReadJSON(&req)
	if err != nil {
		return err
	}
	err = c.ClusterNodeService.Batch(clusterName, req)
	if err != nil {
		return err
	}
	operator := c.Ctx.Values().GetString("operator")
	if req.Operation == "delete" {
		go kolog.Save(operator, constant.DELETE_CLUSTER_NODE, clusterName)
	} else {
		go kolog.Save(operator, constant.CREATE_CLUSTER_NODE, clusterName)
	}

	return nil
}

func (c ClusterController) GetWebkubectlBy(clusterName string) (*dto.WebkubectlToken, error) {
	tk, err := c.ClusterService.GetWebkubectlToken(clusterName)
	if err != nil {
		return nil, err
	}

	return &tk, nil
}

func (c ClusterController) GetSecretBy(clusterName string) (*dto.ClusterSecret, error) {
	sec, err := c.ClusterService.GetSecrets(clusterName)
	if err != nil {
		return nil, err
	}
	return &sec, nil
}

func (c ClusterController) GetCisBy(clusterName string) (*page.Page, error) {
	p, _ := c.Ctx.Values().GetBool("page")
	if p {
		num, _ := c.Ctx.Values().GetInt(constant.PageNumQueryKey)
		size, _ := c.Ctx.Values().GetInt(constant.PageSizeQueryKey)
		pageItem, err := c.CisService.Page(num, size, clusterName)
		if err != nil {
			return nil, err
		}
		return pageItem, nil
	} else {
		var pageItem page.Page
		items, err := c.CisService.List(clusterName)
		if err != nil {
			return nil, err
		}
		pageItem.Items = items
		pageItem.Total = len(items)
		return &pageItem, nil
	}
}

func (c ClusterController) DeleteCisBy(clusterName string, id string) error {
	if clusterName == "" || id == "" {
		return errors.New("params is not set")
	}

	operator := c.Ctx.Values().GetString("operator")
	go kolog.Save(operator, constant.DELETE_CLUSTER_CIS_SCAN_RESULT, clusterName+"-"+id)

	return c.CisService.Delete(clusterName, id)
}

func (c ClusterController) GetCisdetailBy(clusterName string, id string) (*dto.CisTaskDetail, error) {
	if clusterName == "" || id == "" {
		return nil, errors.New("params is not set")
	}
	return c.CisService.Get(clusterName, id)
}

func (c ClusterController) GetCisreportBy(clusterName, id string) error {
	format := "json"
	if c.Ctx.URLParamExists("format") {
		format = c.Ctx.URLParam("format")
	}
	var buf []byte
	var err error
	var t *dto.CisTaskDetail
	t, err = c.CisService.Get(clusterName, id)
	if err != nil {
		return err
	}
	fileName := fmt.Sprintf("cis_report_%s_%s_%s.%s", t.ClusterName, t.Policy, t.ID, format)
	c.Ctx.Header("Content-Type", "application/octet-stream")
	c.Ctx.Header("Content-Disposition", "attachment; filename=\""+fileName+"\"")

	switch format {
	case "json":
		buf, err = json.Marshal(t.CisReport)
		if err != nil {
			return err
		}
	case "yaml", "yml":
		buf, err = yaml.Marshal(t.CisReport)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("can not support formater %s", format)
	}
	_, _ = c.Ctx.Write(buf)
	return nil
}

func (c ClusterController) PostCisBy(clusterName string) (*dto.CisTask, error) {
	var req dto.CisTaskCreate
	if err := c.Ctx.ReadJSON(&req); err != nil {
		return nil, err
	}

	operator := c.Ctx.Values().GetString("operator")
	go kolog.Save(operator, constant.START_CLUSTER_CIS_SCAN, clusterName)

	return c.CisService.Create(clusterName, &req)
}

func (c *ClusterController) GetHealthBy(clusterName string) (*dto.ClusterHealth, error) {
	operator := c.Ctx.Values().GetString("operator")
	go kolog.Save(operator, constant.HEALTH_CHECK, clusterName)

	return c.ClusterHealthService.HealthCheck(clusterName)
}

func (c *ClusterController) PostRecoverBy(clusterName string) ([]dto.ClusterRecoverItem, error) {
	var req dto.ClusterHealth
	err := c.Ctx.ReadJSON(&req)
	if err != nil {
		return []dto.ClusterRecoverItem{}, err
	}
	operator := c.Ctx.Values().GetString("operator")
	go kolog.Save(operator, constant.HEALTH_RECOVER, clusterName)

	return c.ClusterHealthService.Recover(clusterName, req)
}

func (c *ClusterController) GetBackupaccountsBy(name string) ([]dto.BackupAccount, error) {
	return c.BackupAccountService.ListByClusterName(name)
}
