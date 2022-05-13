package controller

import (
	"fmt"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/controller/kolog"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/logger"
	"github.com/KubeOperator/KubeOperator/pkg/service"
	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12/context"
)

type ProvisionerController struct {
	Ctx                              context.Context
	ClusterStorageProvisionerService service.ClusterStorageProvisionerService
}

func NewProvisionerController() *ProvisionerController {
	return &ProvisionerController{
		ClusterStorageProvisionerService: service.NewClusterStorageProvisionerService(),
	}
}

func (c ProvisionerController) GetBy(name string) ([]dto.ClusterStorageProvisioner, error) {
	csp, err := c.ClusterStorageProvisionerService.ListStorageProvisioner(name)
	if err != nil {
		logger.Log.Info(fmt.Sprintf("%+v", err))
		return nil, err
	}
	return csp, nil
}

func (c ProvisionerController) PostBy(name string) (*dto.ClusterStorageProvisioner, error) {
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

func (c ProvisionerController) PostSyncBy(name string) error {
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

func (c ProvisionerController) PostDeleteBy(clusterName string) error {
	var item dto.ClusterStorageProvisioner
	if err := c.Ctx.ReadJSON(&item); err != nil {
		logger.Log.Info(fmt.Sprintf("%+v", err))
		return err
	}
	operator := c.Ctx.Values().GetString("operator")
	go kolog.Save(operator, constant.DELETE_CLUSTER_STORAGE_SUPPLIER, clusterName+"-"+item.Name)

	return c.ClusterStorageProvisionerService.DeleteStorageProvisioner(clusterName, item.Name)
}

func (c ProvisionerController) PostBatchBy(clusterName string) error {
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

func (c ProvisionerController) PostDeployment() (interface{}, error) {
	var req dto.DeploymentSearch
	if err := c.Ctx.ReadJSON(&req); err != nil {
		return nil, err
	}
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return nil, err
	}

	return c.ClusterStorageProvisionerService.SearchDeployment(req)
}
