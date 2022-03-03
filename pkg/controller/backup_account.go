package controller

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/controller/kolog"
	"github.com/KubeOperator/KubeOperator/pkg/controller/koregexp"
	"github.com/KubeOperator/KubeOperator/pkg/controller/page"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/service"
	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12/context"
)

type BackupAccountController struct {
	Ctx                  context.Context
	BackupAccountService service.BackupAccountService
}

func NewBackupAccountController() *BackupAccountController {
	return &BackupAccountController{
		BackupAccountService: service.NewBackupAccountService(),
	}
}

// List BackupAccount
// @Tags backupAccounts
// @Summary Show all backupAccounts
// @Description Show backupAccounts
// @Accept  json
// @Produce  json
// @Param  pageNum  query  int  true "page number"
// @Param  pageSize  query  int  true "page size"
// @Success 200 {object} page.Page
// @Security ApiKeyAuth
// @Router /backupaccounts [get]
func (b BackupAccountController) Get() (page.Page, error) {
	pg, _ := b.Ctx.Values().GetBool("page")
	if pg {
		num, _ := b.Ctx.Values().GetInt(constant.PageNumQueryKey)
		size, _ := b.Ctx.Values().GetInt(constant.PageSizeQueryKey)
		return b.BackupAccountService.Page(num, size)
	} else {
		var page page.Page
		projectName := b.Ctx.URLParam("projectName")
		items, err := b.BackupAccountService.List(projectName)
		if err != nil {
			return page, err
		}
		page.Items = items
		page.Total = len(items)
		return page, nil
	}
}

// Create BackupAccount
// @Tags backupAccounts
// @Summary Create a backupAccount
// @Description create a backupAccount
// @Accept  json
// @Produce  json
// @Param request body dto.BackupAccountCreate true "request"
// @Success 200 {object} dto.BackupAccount
// @Security ApiKeyAuth
// @Router /backupaccounts/ [post]
func (b BackupAccountController) Post() (*dto.BackupAccount, error) {
	var req dto.BackupAccountCreate
	err := b.Ctx.ReadJSON(&req)
	if err != nil {
		return nil, err
	}
	validate := validator.New()
	if err := validate.RegisterValidation("koname", koregexp.CheckNamePattern); err != nil {
		return nil, err
	}
	if err := validate.Struct(req); err != nil {
		return nil, err
	}

	go kolog.Save(b.Ctx, constant.CREATE_BACKUP_ACCOUNT, req.Name)

	return b.BackupAccountService.Create(req)
}

// Delete BackupAccounts
// @Tags backupAccounts
// @Summary Delete backupAccount list
// @Description delete  backupAccount list
// @Accept  json
// @Produce  json
// @Param request body dto.BackupAccountOp true "request"
// @Security ApiKeyAuth
// @Router /backupaccounts/batch [post]
func (b BackupAccountController) PostBatch() error {
	var req dto.BackupAccountOp
	err := b.Ctx.ReadJSON(&req)
	if err != nil {
		return err
	}
	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		return err
	}
	err = b.BackupAccountService.Batch(req)
	if err != nil {
		return err
	}

	delAccs := ""
	for _, item := range req.Items {
		delAccs += (item.Name + ",")
	}
	go kolog.Save(b.Ctx, constant.DELETE_BACKUP_ACCOUNT, delAccs)

	return err
}

// Update BackupAccount
// @Tags backupAccounts
// @Summary Update a backupAccount
// @Description Update a backupAccount
// @Accept  json
// @Produce  json
// @Param request body dto.BackupAccountUpdate true "request"
// @Success 200 {object} dto.BackupAccount
// @Security ApiKeyAuth
// @Router /backupaccounts/{name}/ [patch]
func (b BackupAccountController) PatchBy(name string) (*dto.BackupAccount, error) {
	var req dto.BackupAccountUpdate
	err := b.Ctx.ReadJSON(&req)
	if err != nil {
		return nil, err
	}
	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		return nil, err
	}

	go kolog.Save(b.Ctx, constant.UPDATE_BACKUP_ACCOUNT, name)

	return b.BackupAccountService.Update(req)
}

func (b BackupAccountController) Delete(name string) error {
	go kolog.Save(b.Ctx, constant.DELETE_BACKUP_ACCOUNT, name)

	return b.BackupAccountService.Delete(name)
}

func (b BackupAccountController) PostBuckets() ([]interface{}, error) {
	var req dto.CloudStorageRequest
	err := b.Ctx.ReadJSON(&req)
	if err != nil {
		return nil, err
	}
	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		return nil, err
	}
	return b.BackupAccountService.GetBuckets(req)
}
