package controller

import (
	"errors"
	"fmt"
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/controller/page"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/service"
	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12/context"
)

var (
	HostAlreadyExistsErr = "HOST_ALREADY_EXISTS"
)

type HostController struct {
	Ctx                  context.Context
	HostService          service.HostService
	SystemSettingService service.SystemSettingService
}

func NewHostController() *HostController {
	return &HostController{
		HostService:          service.NewHostService(),
		SystemSettingService: service.NewSystemSettingService(),
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
		projectName := h.Ctx.URLParam("projectName")
		items, err := h.HostService.List(projectName)
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

func (h HostController) Post() (*dto.Host, error) {
	var req dto.HostCreate
	err := h.Ctx.ReadJSON(&req)
	if err != nil {
		return nil, err
	}
	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		return nil, err
	}

	localIp, err := h.SystemSettingService.Get("ip")
	if err != nil {
		return nil, err
	}
	if localIp.Value == req.Ip {
		return nil, errors.New(fmt.Sprintf("%s is localIp, can not imported", localIp))
	}
	item, _ := h.HostService.Get(req.Name)
	if item.ID != "" {
		return nil, errors.New(HostAlreadyExistsErr)
	}
	item, err = h.HostService.Create(req)
	if err != nil {
		return nil, err
	}
	return &item, nil
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
		return err
	}
	return err
}
