package controller

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
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

func (b BackupAccountController) Post() (*dto.BackupAccount, error) {
	var req dto.BackupAccountCreate
	err := b.Ctx.ReadJSON(&req)
	if err != nil {
		return nil, err
	}
	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		return nil, err
	}
	return b.BackupAccountService.Create(req)
}

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
	return err
}
