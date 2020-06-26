package controller

import (
	"errors"
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/controller/page"
	"github.com/KubeOperator/KubeOperator/pkg/controller/warp"
	"github.com/KubeOperator/KubeOperator/pkg/service"
	"github.com/KubeOperator/KubeOperator/pkg/service/dto"
	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12/context"
)

type RegionController struct {
	Ctx           context.Context
	RegionService service.RegionService
}

func NewRegionController() *RegionController {
	return &RegionController{
		RegionService: service.NewRegionService(),
	}
}

func (r RegionController) Get() (page.Page, error) {

	p, _ := r.Ctx.Values().GetBool("page")
	if p {
		num, _ := r.Ctx.Values().GetInt(constant.PageNumQueryKey)
		size, _ := r.Ctx.Values().GetInt(constant.PageSizeQueryKey)
		return r.RegionService.Page(num, size)
	} else {
		var page page.Page
		items, err := r.RegionService.List()
		if err != nil {
			return page, err
		}
		page.Items = items
		page.Total = len(items)
		return page, nil
	}
}

func (r RegionController) GetBy(name string) (dto.Region, error) {
	return r.RegionService.Get(name)
}

func (r RegionController) Post() (dto.Region, error) {
	var req dto.RegionCreate
	err := r.Ctx.ReadJSON(&req)
	if err != nil {
		return dto.Region{}, err
	}
	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		return dto.Region{}, err
	}
	return r.RegionService.Create(req)
}

func (r RegionController) Delete(name string) error {
	return r.RegionService.Delete(name)
}

func (r RegionController) PostBatch() error {
	var req dto.RegionOp
	err := r.Ctx.ReadJSON(&req)
	if err != nil {
		return err
	}
	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		return err
	}
	err = r.RegionService.Batch(req)
	if err != nil {
		return warp.NewControllerError(errors.New(r.Ctx.Tr(err.Error())))
	}
	return err
}

func (r RegionController) PostCheckValid() error {
	var req dto.RegionCreate
	err := r.Ctx.ReadJSON(&req)
	if err != nil {
		return err
	}

	err = r.RegionService.CheckValid(req)
	if err != nil {
		return err
	}
	return nil
}
