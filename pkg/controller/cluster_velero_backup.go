package controller

import (
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/service"
	"github.com/kataras/iris/v12/context"
)

type ClusterVeleroBackupController struct {
	Ctx                 context.Context
	VeleroBackupService service.VeleroBackupService
}

func NewClusterVeleroBackupController() *ClusterVeleroBackupController {
	return &ClusterVeleroBackupController{
		VeleroBackupService: service.NewVeleroBackupService(),
	}
}

func (c ClusterVeleroBackupController) Get() (interface{}, error) {
	clusterName := c.Ctx.Params().GetString("cluster")
	operate := c.Ctx.Params().GetString("operate")
	return c.VeleroBackupService.Get(clusterName, operate)
}

func (c ClusterVeleroBackupController) GetLogs() (string, error) {
	clusterName := c.Ctx.Params().GetString("cluster")
	operate := c.Ctx.Params().GetString("operate")
	name := c.Ctx.URLParam("name")
	return c.VeleroBackupService.GetLogs(clusterName, name, operate)
}

func (c ClusterVeleroBackupController) GetDescribe() (string, error) {
	clusterName := c.Ctx.Params().GetString("cluster")
	operate := c.Ctx.Params().GetString("operate")
	name := c.Ctx.URLParam("name")
	return c.VeleroBackupService.GetDescribe(clusterName, name, operate)
}

func (c ClusterVeleroBackupController) DeleteDel() (string, error) {
	clusterName := c.Ctx.Params().GetString("cluster")
	operate := c.Ctx.Params().GetString("operate")
	name := c.Ctx.URLParam("name")
	return c.VeleroBackupService.Delete(clusterName, name, operate)
}

func (c ClusterVeleroBackupController) PostInstallConfig() (string, error) {
	clusterName := c.Ctx.Params().GetString("cluster")
	var req dto.VeleroInstall
	err := c.Ctx.ReadJSON(&req)
	if err != nil {
		return "", err
	}
	return c.VeleroBackupService.Install(clusterName, req)
}

func (c ClusterVeleroBackupController) GetInstallConfig() (model.ClusterVelero, error) {
	clusterName := c.Ctx.Params().GetString("cluster")
	return c.VeleroBackupService.GetConfig(clusterName)
}

func (c ClusterVeleroBackupController) PostCreate() (string, error) {
	operate := c.Ctx.Params().GetString("operate")
	var req dto.VeleroBackup
	err := c.Ctx.ReadJSON(&req)
	if err != nil {
		return "", err
	}
	clusterName := c.Ctx.Params().GetString("cluster")
	req.Cluster = clusterName
	return c.VeleroBackupService.Create(operate, req)
}
