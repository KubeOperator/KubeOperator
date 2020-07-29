package controller

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/service"
	"github.com/kataras/iris/v12/context"
)

type ClusterController struct {
	Ctx                              context.Context
	ClusterService                   service.ClusterService
	ClusterInitService               service.ClusterInitService
	ClusterStorageProvisionerService service.ClusterStorageProvisionerService
	ClusterToolService               service.ClusterToolService
	ClusterNodeService               service.ClusterNodeService
}

func NewClusterController() *ClusterController {
	return &ClusterController{
		ClusterService:                   service.NewClusterService(),
		ClusterInitService:               service.NewClusterInitService(),
		ClusterStorageProvisionerService: service.NewClusterStorageProvisionerService(),
		ClusterToolService:               service.NewClusterToolService(),
		ClusterNodeService:               service.NewClusterNodeService(),
	}
}

func (c ClusterController) Get() (dto.ClusterPage, error) {
	page, _ := c.Ctx.Values().GetBool("page")
	if page {
		projectName := c.Ctx.URLParam("projectName")
		num, _ := c.Ctx.Values().GetInt(constant.PageNumQueryKey)
		size, _ := c.Ctx.Values().GetInt(constant.PageSizeQueryKey)
		return c.ClusterService.Page(num, size, projectName)
	} else {
		var page dto.ClusterPage
		items, err := c.ClusterService.List()
		if err != nil {
			return page, err
		}
		page.Items = items
		page.Total = len(items)
		return page, nil
	}
}

func (c ClusterController) GetBy(name string) (dto.Cluster, error) {
	return c.ClusterService.Get(name)
}

func (c ClusterController) GetStatusBy(name string) (dto.ClusterStatus, error) {
	return c.ClusterService.GetStatus(name)
}

func (c ClusterController) Post() error {
	var req dto.ClusterCreate
	err := c.Ctx.ReadJSON(&req)
	if err != nil {
		return err
	}
	return c.ClusterService.Create(req)
}

func (c ClusterController) PostInitBy(name string) error {
	return c.ClusterInitService.Init(name)
}

func (c ClusterController) GetProvisionerBy(name string) ([]dto.ClusterStorageProvisioner, error) {
	return c.ClusterStorageProvisionerService.ListStorageProvisioner(name)
}
func (c ClusterController) PostProvisionerBy(name string) (dto.ClusterStorageProvisioner, error) {
	var req dto.ClusterStorageProvisionerCreation
	err := c.Ctx.ReadJSON(&req)
	if err != nil {
		return dto.ClusterStorageProvisioner{}, err
	}
	return c.ClusterStorageProvisionerService.CreateStorageProvisioner(name, req)
}

func (c ClusterController) DeleteProvisionerBy(clusterName string, name string) error {
	return c.ClusterStorageProvisionerService.DeleteStorageProvisioner(clusterName, name)
}

func (c ClusterController) PostProvisionerBatchBy(clusterName string) error {
	var batch dto.ClusterStorageProvisionerBatch
	if err := c.Ctx.ReadJSON(&batch); err != nil {
		return err
	}
	return c.ClusterStorageProvisionerService.BatchStorageProvisioner(clusterName, batch)
}

func (c ClusterController) GetToolBy(clusterName string) ([]dto.ClusterTool, error) {
	return c.ClusterToolService.List(clusterName)
}

func (c ClusterController) PostToolBy(clusterName string) (dto.ClusterTool, error) {
	var req dto.ClusterTool
	if err := c.Ctx.ReadJSON(&req); err != nil {
		return req, err
	}
	return c.ClusterToolService.Enable(clusterName, req)
}

func (c ClusterController) Delete(name string) error {
	return c.ClusterService.Delete(name)
}

func (c ClusterController) PostBatch() error {
	var batch dto.ClusterBatch
	if err := c.Ctx.ReadJSON(&batch); err != nil {
		return err
	}
	if err := c.ClusterService.Batch(batch); err != nil {
		return err
	}
	return nil
}

func (c ClusterController) GetNodeBy(clusterName string) ([]dto.Node, error) {
	return c.ClusterNodeService.List(clusterName)
}

func (c ClusterController) PostNodeBatchBy(clusterName string) ([]dto.Node, error) {
	var req dto.NodeBatch
	err := c.Ctx.ReadJSON(&req)
	if err != nil {
		return nil, err
	}
	return c.ClusterNodeService.Batch(clusterName, req)
}

func (c ClusterController) GetWebkubectlBy(clusterName string) (dto.WebkubectlToken, error) {
	return c.ClusterService.GetWebkubectlToken(clusterName)
}

func (c ClusterController) GetSecretBy(clusterName string) (dto.ClusterSecret, error) {
	return c.ClusterService.GetSecrets(clusterName)
}
