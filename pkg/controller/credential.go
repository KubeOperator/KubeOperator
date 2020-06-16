package controller

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/service"
	"github.com/KubeOperator/KubeOperator/pkg/service/dto"
	"github.com/kataras/iris/v12/context"
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

func (c CredentialController) Get() (dto.CredentialPage, error) {

	page, _ := c.Ctx.Values().GetBool("page")
	if page {
		num, _ := c.Ctx.Values().GetInt(constant.PageNumQueryKey)
		size, _ := c.Ctx.Values().GetInt(constant.PageSizeQueryKey)
		return c.CredentialService.Page(num, size)
	} else {
		var page dto.CredentialPage
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

func (c CredentialController) Post() error {
	var req dto.CredentialCreate
	err := c.Ctx.ReadJSON(&req)
	if err != nil {
		return err
	}
	return c.CredentialService.Create(req)
}

func (c CredentialController) Delete(name string) error {
	return c.CredentialService.Delete(name)
}

func (c CredentialController) Patch() error {
	var req dto.CredentialUpdate
	err := c.Ctx.ReadJSON(&req)
	if err != nil {
		return err
	}
	return c.CredentialService.Update(req)
}

func (c CredentialController) Batch(operation string, items []dto.Credential) error {
	_, err := c.CredentialService.Batch(operation, items)
	return err
}
