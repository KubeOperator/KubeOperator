package controller

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/controller/page"
	"github.com/KubeOperator/KubeOperator/pkg/controller/warp"
	"github.com/KubeOperator/KubeOperator/pkg/service"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12/context"
	"golang.org/x/tools/go/ssa/interp/testdata/src/errors"
)

type CredentialController struct {
	Ctx               context.Context
	CredentialService service.CredentialService
}

func NewCredentialController() *CredentialController {
	return &CredentialController{
		CredentialService: service.NewCredentialService(),
	}
}

func (c CredentialController) Get() (page.Page, error) {

	p, _ := c.Ctx.Values().GetBool("page")
	if p {
		num, _ := c.Ctx.Values().GetInt(constant.PageNumQueryKey)
		size, _ := c.Ctx.Values().GetInt(constant.PageSizeQueryKey)
		return c.CredentialService.Page(num, size)
	} else {
		var page page.Page
		items, err := c.CredentialService.List()
		if err != nil {
			return page, err
		}
		page.Items = items
		page.Total = len(items)
		return page, nil
	}
}

func (c CredentialController) GetBy(name string) (dto.Credential, error) {
	return c.CredentialService.Get(name)
}

func (c CredentialController) Post() (dto.Credential, error) {
	var req dto.CredentialCreate
	err := c.Ctx.ReadJSON(&req)
	if err != nil {
		return dto.Credential{}, err
	}
	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		return dto.Credential{}, err
	}
	return c.CredentialService.Create(req)
}

func (c CredentialController) Delete(name string) error {
	return c.CredentialService.Delete(name)
}

func (c CredentialController) PatchBy(name string) (dto.Credential, error) {
	var req dto.CredentialUpdate
	err := c.Ctx.ReadJSON(&req)
	if err != nil {
		return dto.Credential{}, err
	}
	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		return dto.Credential{}, err
	}
	return c.CredentialService.Update(req)
}

func (c CredentialController) PostBatch() error {
	var req dto.CredentialBatchOp
	err := c.Ctx.ReadJSON(&req)
	if err != nil {
		return err
	}
	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		return err
	}
	err = c.CredentialService.Batch(req)
	if err != nil {
		return warp.NewControllerError(errors.New(c.Ctx.Tr(err.Error())))
	}
	return err
}
