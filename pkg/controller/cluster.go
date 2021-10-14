package controller

import (
	"errors"
	"fmt"
	"io"

	"github.com/KubeOperator/KubeOperator/pkg/controller/condition"
	"github.com/KubeOperator/KubeOperator/pkg/logger"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/controller/kolog"
	"github.com/KubeOperator/KubeOperator/pkg/controller/page"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/service"
	"github.com/KubeOperator/KubeOperator/pkg/util/ansible"
	"github.com/kataras/iris/v12/context"
)

type ClusterController struct {
	Ctx                              context.Context
	ClusterService                   service.ClusterService
	ClusterInitService               service.ClusterInitService
	ClusterStorageProvisionerService service.ClusterStorageProvisionerService
	ClusterToolService               service.ClusterToolService
	ClusterNodeService               service.ClusterNodeService
	ClusterLogService                service.ClusterLogService
	ClusterImportService             service.ClusterImportService
	CisService                       service.CisService
	ClusterUpgradeService            service.ClusterUpgradeService
	ClusterHealthService             service.ClusterHealthService
	BackupAccountService             service.BackupAccountService
}

func NewClusterController() *ClusterController {
	return &ClusterController{
		ClusterService:                   service.NewClusterService(),
		ClusterInitService:               service.NewClusterInitService(),
		ClusterStorageProvisionerService: service.NewClusterStorageProvisionerService(),
		ClusterToolService:               service.NewClusterToolService(),
		ClusterNodeService:               service.NewClusterNodeService(),
		ClusterLogService:                service.NewClusterLogService(),
		ClusterImportService:             service.NewClusterImportService(),
		CisService:                       service.NewCisService(),
		ClusterUpgradeService:            service.NewClusterUpgradeService(),
		ClusterHealthService:             service.NewClusterHealthService(),
		BackupAccountService:             service.NewBackupAccountService(),
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

func (c ClusterController) GetStatusBy(name string) (*dto.ClusterStatus, error) {
	cs, err := c.ClusterService.GetStatus(name)
	if err != nil {
		return nil, err
	}
	return &cs, nil
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

	return c.ClusterInitService.Init(name)
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

func (c ClusterController) GetProvisionerBy(name string) ([]dto.ClusterStorageProvisioner, error) {
	csp, err := c.ClusterStorageProvisionerService.ListStorageProvisioner(name)
	if err != nil {
		logger.Log.Info(fmt.Sprintf("%+v", err))
		return nil, err
	}
	return csp, nil
}
func (c ClusterController) PostProvisionerBy(name string) (*dto.ClusterStorageProvisioner, error) {
	var req dto.ClusterStorageProvisionerCreation
	err := c.Ctx.ReadJSON(&req)
	if err != nil {
		return nil, err
	}
	p, err := c.ClusterStorageProvisionerService.CreateStorageProvisioner(name, req)
	if err != nil {
		logger.Log.Info(fmt.Sprintf("%+v", err))
		return nil, err
	}

	operator := c.Ctx.Values().GetString("operator")
	go kolog.Save(operator, constant.CREATE_CLUSTER_STORAGE_SUPPLIER, name+"-"+req.Name+"("+req.Type+")")

	return &p, nil
}
func (c ClusterController) PostProvisionerSyncBy(name string) error {
	var req []dto.ClusterStorageProvisionerSync
	err := c.Ctx.ReadJSON(&req)
	if err != nil {
		return err
	}
	if err := c.ClusterStorageProvisionerService.SyncStorageProvisioner(name, req); err != nil {
		logger.Log.Info(fmt.Sprintf("%+v", err))
		return err
	}

	var proStr string
	for _, pro := range req {
		proStr += (pro.Name + ",")
	}
	operator := c.Ctx.Values().GetString("operator")
	go kolog.Save(operator, constant.SYNC_CLUSTER_STORAGE_SUPPLIER, proStr)

	return nil
}

func (c ClusterController) PostProvisionerDeleteBy(clusterName string) error {
	var item dto.ClusterStorageProvisioner
	if err := c.Ctx.ReadJSON(&item); err != nil {
		logger.Log.Info(fmt.Sprintf("%+v", err))
		return err
	}
	operator := c.Ctx.Values().GetString("operator")
	go kolog.Save(operator, constant.DELETE_CLUSTER_STORAGE_SUPPLIER, clusterName+"-"+item.Name)

	return c.ClusterStorageProvisionerService.DeleteStorageProvisioner(clusterName, item.Name)
}

func (c ClusterController) PostProvisionerBatchBy(clusterName string) error {
	var batch dto.ClusterStorageProvisionerBatch
	if err := c.Ctx.ReadJSON(&batch); err != nil {
		return err
	}
	if err := c.ClusterStorageProvisionerService.BatchStorageProvisioner(clusterName, batch); err != nil {
		logger.Log.Info(fmt.Sprintf("%+v", err))
		return err
	}

	operator := c.Ctx.Values().GetString("operator")
	delClus := ""
	for _, item := range batch.Items {
		delClus += (item.Name + ",")
	}
	go kolog.Save(operator, constant.DELETE_CLUSTER_STORAGE_SUPPLIER, clusterName+"-"+delClus)

	return nil
}

func (c ClusterController) GetToolBy(clusterName string) ([]dto.ClusterTool, error) {
	cts, err := c.ClusterToolService.List(clusterName)
	if err != nil {
		logger.Log.Info(fmt.Sprintf("%+v", err))
		return nil, err
	}
	return cts, nil
}

func (c ClusterController) GetToolPortBy(clusterName, namespace, toolName, toolVersion string) (dto.ClusterTool, error) {
	tool, err := c.ClusterToolService.GetNodePort(clusterName, toolName, toolVersion, namespace)
	if err != nil {
		logger.Log.Info(fmt.Sprintf("%+v", err))
		return tool, err
	}
	return tool, nil
}

func (c ClusterController) PostToolSyncBy(clusterName string) (*[]dto.ClusterTool, error) {
	cts, err := c.ClusterToolService.SyncStatus(clusterName)
	if err != nil {
		logger.Log.Info(fmt.Sprintf("%+v", err))
		return nil, err
	}

	return &cts, nil
}

func (c ClusterController) PostToolEnableBy(clusterName string) (*dto.ClusterTool, error) {
	var req dto.ClusterTool
	if err := c.Ctx.ReadJSON(&req); err != nil {
		return nil, err
	}
	cts, err := c.ClusterToolService.Enable(clusterName, req)
	if err != nil {
		logger.Log.Info(fmt.Sprintf("%+v", err))
		return nil, err
	}

	operator := c.Ctx.Values().GetString("operator")
	go kolog.Save(operator, constant.ENABLE_CLUSTER_TOOL, clusterName+"-"+req.Name)

	return &cts, nil
}

func (c ClusterController) PostToolUpgradeBy(clusterName string) (*dto.ClusterTool, error) {
	var req dto.ClusterTool
	if err := c.Ctx.ReadJSON(&req); err != nil {
		return nil, err
	}
	cts, err := c.ClusterToolService.Upgrade(clusterName, req)
	if err != nil {
		logger.Log.Info(fmt.Sprintf("%+v", err))
		return nil, err
	}

	operator := c.Ctx.Values().GetString("operator")
	go kolog.Save(operator, constant.UPGRADE_CLUSTER_TOOL, clusterName+"-"+req.Name)

	return &cts, nil
}

func (c ClusterController) PostToolDisableBy(clusterName string) (*dto.ClusterTool, error) {
	var req dto.ClusterTool
	if err := c.Ctx.ReadJSON(&req); err != nil {
		return nil, err
	}
	cts, err := c.ClusterToolService.Disable(clusterName, req)
	if err != nil {
		return nil, err
	}

	operator := c.Ctx.Values().GetString("operator")
	go kolog.Save(operator, constant.DISABLE_CLUSTER_TOOL, clusterName+"-"+req.Name)

	return &cts, nil
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
// @Security ApiKeyAuth
// @Success 200
// @Router /clusters/node/recreate/{clusterName}/{node} [post]
func (c ClusterController) PostNodeRecreateBy(clusterName string, node string) error {
	if err := c.ClusterNodeService.Recreate(clusterName, node); err != nil {
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

func (c ClusterController) GetLogBy(clusterName string) ([]dto.ClusterLog, error) {
	ls, err := c.ClusterLogService.List(clusterName)
	if err != nil {
		return nil, err
	}
	return ls, nil

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

func (c ClusterController) PostCisBy(clusterName string) (*dto.CisTask, error) {
	operator := c.Ctx.Values().GetString("operator")
	go kolog.Save(operator, constant.START_CLUSTER_CIS_SCAN, clusterName)

	return c.CisService.Create(clusterName)
}

type Log struct {
	Msg string `json:"msg"`
}

func (c ClusterController) GetLoggerBy(clusterName string) (*Log, error) {
	cluster, err := c.ClusterService.Get(clusterName)
	if err != nil {
		return nil, err
	}
	r, err := ansible.GetAnsibleLogReader(cluster.Name, cluster.LogId)
	if err != nil {
		return nil, err
	}
	var chunk []byte
	for {

		buffer := make([]byte, 1024)
		n, err := r.Read(buffer)
		if err != nil && err != io.EOF {
			return nil, err
		}
		if n == 0 {
			break
		}
		chunk = append(chunk, buffer[:n]...)
	}
	return &Log{Msg: string(chunk)}, nil
}

func (c ClusterController) GetProvisionerLogBy(clusterName, logId string) (*Log, error) {
	r, err := ansible.GetAnsibleLogReader(clusterName, logId)
	if err != nil {
		return nil, err
	}
	var chunk []byte
	for {

		buffer := make([]byte, 1024)
		n, err := r.Read(buffer)
		if err != nil && err != io.EOF {
			return nil, err
		}
		if n == 0 {
			break
		}
		chunk = append(chunk, buffer[:n]...)
	}
	return &Log{Msg: string(chunk)}, nil
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
