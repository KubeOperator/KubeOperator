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
	ClusterMonitorService            service.ClusterMonitorService
	ClusterStorageProvisionerService service.ClusterStorageProvisionerService
}

func NewClusterController() *ClusterController {
	return &ClusterController{
		ClusterService:                   service.NewClusterService(),
		ClusterInitService:               service.NewClusterInitService(),
		ClusterMonitorService:            service.NewClusterMonitorService(),
		ClusterStorageProvisionerService: service.NewClusterStorageProvisionerService(),
	}
}

func (c ClusterController) Get() (dto.ClusterPage, error) {
	page, _ := c.Ctx.Values().GetBool("page")
	if page {
		num, _ := c.Ctx.Values().GetInt(constant.PageNumQueryKey)
		size, _ := c.Ctx.Values().GetInt(constant.PageSizeQueryKey)
		return c.ClusterService.Page(num, size)
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

func (c ClusterController) GetMonitorBy(name string) (dto.ClusterMonitor, error) {
	return c.ClusterService.GetMonitor(name)
}

func (c ClusterController) PostMonitorBy(name string) error {
	return c.ClusterMonitorService.Init(name)
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
