package controller

import (
	"fmt"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/controller/kolog"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/logger"
	"github.com/KubeOperator/KubeOperator/pkg/service"
	"github.com/kataras/iris/v12/context"
)

type ClusterToolController struct {
	Ctx                context.Context
	ClusterToolService service.ClusterToolService
}

func NewClusterToolController() *ClusterToolController {
	return &ClusterToolController{
		ClusterToolService: service.NewClusterToolService(),
	}
}

func (c ClusterToolController) GetBy(clusterName string) ([]dto.ClusterTool, error) {
	cts, err := c.ClusterToolService.List(clusterName)
	if err != nil {
		logger.Log.Info(fmt.Sprintf("%+v", err))
		return nil, err
	}
	return cts, nil
}

func (c ClusterToolController) GetPortBy(clusterName, toolName string) (string, error) {
	endPoint, err := c.ClusterToolService.GetNodePort(clusterName, toolName)
	if err != nil {
		logger.Log.Info(fmt.Sprintf("%+v", err))
		return endPoint, err
	}
	return endPoint, nil
}

func (c ClusterToolController) PostSyncBy(clusterName string) (*[]dto.ClusterTool, error) {
	cts, err := c.ClusterToolService.SyncStatus(clusterName)
	if err != nil {
		logger.Log.Info(fmt.Sprintf("%+v", err))
		return nil, err
	}

	return &cts, nil
}

func (c ClusterToolController) PostEnableBy(clusterName string) (*dto.ClusterTool, error) {
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

func (c ClusterToolController) GetFlexBy(clusterName string) (string, error) {
	return c.ClusterToolService.GetFlex(clusterName)
}

func (c ClusterToolController) PostFlexEnableBy(clusterName string) error {
	if err := c.ClusterToolService.EnableFlex(clusterName); err != nil {
		logger.Log.Info(fmt.Sprintf("%+v", err))
		return err
	}

	operator := c.Ctx.Values().GetString("operator")
	go kolog.Save(operator, constant.ENABLE_CLUSTER_FLEX, clusterName)

	return nil
}

func (c ClusterToolController) PostFlexDisableBy(clusterName string) error {
	if err := c.ClusterToolService.DisableFlex(clusterName); err != nil {
		logger.Log.Info(fmt.Sprintf("%+v", err))
		return err
	}

	operator := c.Ctx.Values().GetString("operator")
	go kolog.Save(operator, constant.DISABLE_CLUSTER_FLEX, clusterName)

	return nil
}

func (c ClusterToolController) PostUpgradeBy(clusterName string) (*dto.ClusterTool, error) {
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

func (c ClusterToolController) PostDisableBy(clusterName string) (*dto.ClusterTool, error) {
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
