package controller

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/controller/kolog"
	"github.com/KubeOperator/KubeOperator/pkg/controller/koregexp"
	"github.com/KubeOperator/KubeOperator/pkg/controller/page"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/service"
	"github.com/KubeOperator/KubeOperator/pkg/util/validator_error"
	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12/context"
)

type VmConfigController struct {
	Ctx             context.Context
	VmConfigService service.VmConfigService
}

func NewVmConfigController() *VmConfigController {
	return &VmConfigController{
		VmConfigService: service.NewVmConfigService(),
	}
}

func (v VmConfigController) Get() (page.Page, error) {
	pa, _ := v.Ctx.Values().GetBool("page")
	if pa {
		num, _ := v.Ctx.Values().GetInt(constant.PageNumQueryKey)
		size, _ := v.Ctx.Values().GetInt(constant.PageSizeQueryKey)
		return v.VmConfigService.Page(num, size)
	} else {
		var page page.Page
		items, err := v.VmConfigService.List()
		if err != nil {
			return page, err
		}
		page.Items = items
		page.Total = len(items)
		return page, nil
	}
}

func (v VmConfigController) Post() (*dto.VmConfig, error) {
	var req dto.VmConfigCreate
	err := v.Ctx.ReadJSON(&req)
	if err != nil {
		return nil, err
	}
	validate := validator.New()
	if err := validate.RegisterValidation("kovmconfig", koregexp.CheckVmConfigPattern); err != nil {
		return nil, err
	}
	if err := validate.Struct(req); err != nil {
		return nil, err
	}

	go kolog.Save(v.Ctx, constant.CREATE_VM_CONFIG, req.Name)

	return v.VmConfigService.Create(req)
}

func (v VmConfigController) PatchBy(name string) (*dto.VmConfig, error) {
	var req dto.VmConfigUpdate
	err := v.Ctx.ReadJSON(&req)
	if err != nil {
		return nil, err
	}
	validate := validator.New()
	validator_error.RegisterTagNameFunc(v.Ctx, validate)
	err = validate.Struct(req)
	if err != nil {
		return nil, validator_error.Tr(v.Ctx, validate, err)
	}
	result, err := v.VmConfigService.Update(req)
	if err != nil {
		return nil, err
	}

	go kolog.Save(v.Ctx, constant.UPDATE_VM_CONFIG, name)

	return result, nil
}

func (v VmConfigController) PostBatch() error {
	var req dto.VmConfigOp
	err := v.Ctx.ReadJSON(&req)
	if err != nil {
		return err
	}
	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		return err
	}
	err = v.VmConfigService.Batch(req)
	if err != nil {
		return err
	}

	delConfs := ""
	for _, item := range req.Items {
		delConfs += (item.Name + ",")
	}
	go kolog.Save(v.Ctx, constant.DELETE_VM_CONFIG, delConfs)

	return err
}
