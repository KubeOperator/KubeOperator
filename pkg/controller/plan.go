package controller

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/controller/page"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/service"
	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12/context"
)

type PlanController struct {
	Ctx         context.Context
	PlanService service.PlanService
}

func NewPlanController() *PlanController {
	return &PlanController{
		PlanService: service.NewPlanService(),
	}
}

func (p PlanController) Get() (page.Page, error) {

	pg, _ := p.Ctx.Values().GetBool("page")
	if pg {
		num, _ := p.Ctx.Values().GetInt(constant.PageNumQueryKey)
		size, _ := p.Ctx.Values().GetInt(constant.PageSizeQueryKey)
		return p.PlanService.Page(num, size)
	} else {
		var page page.Page
		projectName := p.Ctx.URLParam("projectName")
		items, err := p.PlanService.List(projectName)
		if err != nil {
			return page, err
		}
		page.Items = items
		page.Total = len(items)
		return page, nil
	}
}

func (p PlanController) GetBy(name string) (dto.Plan, error) {
	return p.PlanService.Get(name)
}

func (p PlanController) Post() (*dto.Plan, error) {
	var req dto.PlanCreate
	err := p.Ctx.ReadJSON(&req)
	if err != nil {
		return nil, err
	}
	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		return nil, err
	}
	return p.PlanService.Create(req)
}

func (p PlanController) Delete(name string) error {
	return p.PlanService.Delete(name)
}

func (p PlanController) PostBatch() error {
	var req dto.PlanOp
	err := p.Ctx.ReadJSON(&req)
	if err != nil {
		return err
	}
	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		return err
	}
	err = p.PlanService.Batch(req)
	if err != nil {
		return err
	}
	return err
}

func (p PlanController) GetConfigsBy(regionName string) ([]dto.PlanVmConfig, error) {

	return p.PlanService.GetConfigs(regionName)
}
