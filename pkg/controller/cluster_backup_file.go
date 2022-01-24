package controller

import (
	"io/ioutil"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/controller/kolog"
	"github.com/KubeOperator/KubeOperator/pkg/controller/koregexp"
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

// List BackupFile
// @Tags backupFiles
// @Summary Show all backupFiles
// @Description Show backupFiles
// @Accept  json
// @Produce  json
// @Success 200 {object} page.Page
// @Security ApiKeyAuth
// @Router /clusters/backup/files/ [get]
func (b BackupFileController) Get() (*page.Page, error) {
	num, _ := b.Ctx.Values().GetInt(constant.PageNumQueryKey)
	size, _ := b.Ctx.Values().GetInt(constant.PageSizeQueryKey)
	clusterName := b.Ctx.URLParam("clusterName")
	return b.ClusterBackupFileService.Page(num, size, clusterName)
}

// Delete BackupFile
// @Tags backupFiles
// @Summary Delete a BackupFile
// @Description delete a BackupFile by name
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Router /clusters/backup/files/{name}/ [delete]
func (b BackupFileController) DeleteBy(name string) error {
	operator := b.Ctx.Values().GetString("operator")
	go kolog.Save(operator, constant.DELETE_RECOVERY_LIST, name)

	return b.ClusterBackupFileService.Delete(name)
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

	operator := b.Ctx.Values().GetString("operator")
	delBackup := ""
	for _, item := range req.Items {
		delBackup += (item.Name + ",")
	}
	go kolog.Save(operator, constant.DELETE_RECOVERY_LIST, delBackup)

	return err
}

// CLuster Backup
// @Tags backupFiles
// @Summary Backup CLuster
// @Description Backup CLuster
// @Accept  json
// @Produce  json
// @Param request body dto.ClusterBackupFileCreate true "request"
// @Success 200
// @Security ApiKeyAuth
// @Router /clusters/backup/files/backup/ [post]
func (b BackupFileController) PostBackup() error {
	var req dto.ClusterBackupFileCreate
	err := b.Ctx.ReadJSON(&req)
	if err != nil {
		return err
	}
	validate := validator.New()

	if err := validate.RegisterValidation("clustername", koregexp.CheckClusterNamePattern); err != nil {
		return err
	}
	if err := validate.Struct(req); err != nil {
		return err
	}
	err = b.ClusterBackupFileService.Backup(req)
	if err != nil {
		return err
	}

	operator := b.Ctx.Values().GetString("operator")
	if len(req.Name) != 0 {
		go kolog.Save(operator, constant.START_CLUSTER_BACKUP, req.ClusterName+"-"+req.Name)
	} else {
		go kolog.Save(operator, constant.START_CLUSTER_BACKUP, req.ClusterName)
	}

	return err
}

// CLuster Restore
// @Tags backupFiles
// @Summary Restore CLuster
// @Description Restore CLuster
// @Accept  json
// @Produce  json
// @Param request body dto.ClusterBackupFileRestore true "request"
// @Success 200
// @Security ApiKeyAuth
// @Router /clusters/backup/files/restore/ [post]
func (b BackupFileController) PostRestore() error {
	var req dto.ClusterBackupFileRestore
	err := b.Ctx.ReadJSON(&req)
	if err != nil {
		return err
	}
	validate := validator.New()
	if err := validate.RegisterValidation("clustername", koregexp.CheckClusterNamePattern); err != nil {
		return err
	}
	if err := validate.Struct(req); err != nil {
		return err
	}
	err = b.ClusterBackupFileService.Restore(req)
	if err != nil {
		return err
	}

	operator := b.Ctx.Values().GetString("operator")
	go kolog.Save(operator, constant.RECOVER_FROM_RECOVERY, req.ClusterName+"-"+req.Name)

	return err
}

func (b BackupFileController) PostRestoreLocal() error {

	f, _, err := b.Ctx.FormFile("file")
	if err != nil {
		return err
	}
	bs, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	defer f.Close()
	clusterName := b.Ctx.FormValue("clusterName")

	operator := b.Ctx.Values().GetString("operator")
	go kolog.Save(operator, constant.UPLOAD_LOCAL_RECOVERY_FILE, clusterName)

	return b.ClusterBackupFileService.LocalRestore(clusterName, bs)
}
