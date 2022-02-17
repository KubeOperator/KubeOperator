package controller

import (
	"github.com/KubeOperator/KubeOperator/pkg/dto"
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
	clusterName := c.Ctx.Params().GetString("name")
	return c.VeleroBackupService.GetBackups(clusterName)
}

func (c ClusterVeleroBackupController) GetDescribe() (string, error) {
	clusterName := c.Ctx.Params().GetString("name")
	backupName := c.Ctx.URLParam("backupName")
	return c.VeleroBackupService.GetBackupDescribe(clusterName, backupName)
}
func (c ClusterVeleroBackupController) GetLogs() (string, error) {
	clusterName := c.Ctx.Params().GetString("name")
	backupName := c.Ctx.URLParam("backupName")
	return c.VeleroBackupService.GetBackupLogs(clusterName, backupName)
}

func (c ClusterVeleroBackupController) PostCreate() (string, error) {
	var req dto.VeleroBackup
	err := c.Ctx.ReadJSON(&req)
	if err != nil {
		return "", err
	}
	clusterName := c.Ctx.Params().GetString("name")
	req.Cluster = clusterName
	return c.VeleroBackupService.CreateBackup(req)
}
