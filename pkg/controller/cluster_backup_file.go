package controller

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/controller/page"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/service"
	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12/context"
)

type BackupFileController struct {
	Ctx                      context.Context
	ClusterBackupFileService service.CLusterBackupFileService
}

func NewClusterBackupFileController() *BackupFileController {
	return &BackupFileController{
		ClusterBackupFileService: service.NewClusterBackupFileService(),
	}
}

func (b BackupFileController) Get() (*page.Page, error) {
	num, _ := b.Ctx.Values().GetInt(constant.PageNumQueryKey)
	size, _ := b.Ctx.Values().GetInt(constant.PageSizeQueryKey)
	clusterName := b.Ctx.URLParam("clusterName")
	return b.ClusterBackupFileService.Page(num, size, clusterName)
}

func (b BackupFileController) PostBatch() error {
	var req dto.ClusterBackupFileOp
	err := b.Ctx.ReadJSON(&req)
	if err != nil {
		return err
	}
	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		return err
	}
	err = b.ClusterBackupFileService.Batch(req)
	if err != nil {
		return err
	}
	return err
}
