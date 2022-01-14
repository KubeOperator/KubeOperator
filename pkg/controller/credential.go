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

type CredentialController struct {
	Ctx               context.Context
	CredentialService service.CredentialService
}

func NewCredentialController() *CredentialController {
	return &CredentialController{
		CredentialService: service.NewCredentialService(),
	}
}

// List Credential
// @Tags credentials
// @Summary Show all credentials
// @Description Show credentials
// @Accept  json
// @Produce  json
// @Success 200 {object} page.Page
// @Security ApiKeyAuth
// @Router /credentials/ [get]
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

// Create Credential
// @Tags credentials
// @Summary Create a credential
// @Description create a credential
// @Accept  json
// @Produce  json
// @Param request body dto.CredentialCreate true "request"
// @Success 200 {object} dto.Credential
// @Security ApiKeyAuth
// @Router /credentials/ [post]
func (c CredentialController) Post() (*dto.Credential, error) {
	var req dto.CredentialCreate
	err := c.Ctx.ReadJSON(&req)
	if err != nil {
		return nil, err
	}
	validate := validator.New()
	validate.RegisterValidation("koname", koregexp.CheckNamePattern)
	if err := validate.Struct(req); err != nil {
		return nil, err
	}

	operator := c.Ctx.Values().GetString("operator")
	go kolog.Save(operator, constant.CREATE_CREDENTIALS, req.Name)

	return c.CredentialService.Create(req)
}

// Delete Credential
// @Tags credentials
// @Summary Delete a credential
// @Description delete a  credential by name
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Router /backupAccounts/{name}/ [delete]
func (c CredentialController) Delete(name string) error {
	operator := c.Ctx.Values().GetString("operator")
	go kolog.Save(operator, constant.DELETE_CREDENTIALS, name)

	return c.CredentialService.Delete(name)
}

// Update Credential
// @Tags credentials
// @Summary Update a credential
// @Description Update a credential
// @Accept  json
// @Produce  json
// @Param request body dto.CredentialUpdate true "request"
// @Success 200 {object} dto.Credential
// @Security ApiKeyAuth
// @Router /backupAccounts/ [patch]
func (c CredentialController) PatchBy(name string) (dto.Credential, error) {
	var req dto.CredentialUpdate
	err := c.Ctx.ReadJSON(&req)
	if err != nil {
		return dto.Credential{}, err
	}
	validate := validator.New()
	validate.RegisterValidation("koname", koregexp.CheckNamePattern)
	if err := validate.Struct(req); err != nil {
		return dto.Credential{}, err
	}

	operator := c.Ctx.Values().GetString("operator")
	go kolog.Save(operator, constant.UPDATE_CREDENTIALS, name)

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
		return err
	}

	operator := c.Ctx.Values().GetString("operator")
	delCres := ""
	for _, item := range req.Items {
		delCres += (item.Name + ",")
	}
	go kolog.Save(operator, constant.DELETE_CREDENTIALS, delCres)

	return err
}
