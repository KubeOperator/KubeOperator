package controller

import (
	"errors"
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/controller/page"
	"github.com/KubeOperator/KubeOperator/pkg/controller/warp"
	"github.com/KubeOperator/KubeOperator/pkg/service"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12/context"
)

type HostController struct {
	Ctx         context.Context
	HostService service.HostService
}

func NewHostController() *HostController {
	return &HostController{
		HostService: service.NewHostService(),
	}
}

func (h HostController) Get() (page.Page, error) {

	p, _ := h.Ctx.Values().GetBool("page")
	if p {
		num, _ := h.Ctx.Values().GetInt(constant.PageNumQueryKey)
		size, _ := h.Ctx.Values().GetInt(constant.PageSizeQueryKey)
		return h.HostService.Page(num, size)
	} else {
		var page page.Page
		items, err := h.HostService.List()
		if err != nil {
			return page, err
		}
		page.Items = items
		page.Total = len(items)
		return page, nil
	}
}

func (h HostController) GetBy(name string) (dto.Host, error) {
	return h.HostService.Get(name)
}

func (h HostController) Post() (dto.Host, error) {
	var req dto.HostCreate
	err := h.Ctx.ReadJSON(&req)
	if err != nil {
		return dto.Host{}, err
	}
	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		return dto.Host{}, err
	}
	return h.HostService.Create(req)
}

func (h HostController) Delete(name string) error {
	return h.HostService.Delete(name)
}

func (h HostController) PostSyncBy(name string) (dto.Host, error) {
	return h.HostService.Sync(name)
}

func (h HostController) PostBatch() error {
	var req dto.HostOp
	err := h.Ctx.ReadJSON(&req)
	if err != nil {
		return err
	}
	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		return err
	}
	err = h.HostService.Batch(req)
	if err != nil {
		return warp.NewControllerError(errors.New(h.Ctx.Tr(err.Error())))
	}
	return err
}
