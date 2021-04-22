package controller

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/controller/condition"
	"github.com/KubeOperator/KubeOperator/pkg/controller/kolog"
	"github.com/KubeOperator/KubeOperator/pkg/controller/page"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/service"
	sessionUtil "github.com/KubeOperator/KubeOperator/pkg/util/session"
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

// List Plan
// @Tags plans
// @Summary Show all plans
// @Description 获取部署计划列表
// @Accept  json
// @Produce  json
// @Success 200 {object} page.Page
// @Security ApiKeyAuth
// @Router /plans [get]
func (p PlanController) Get() (*page.Page, error) {
	projectName, err := sessionUtil.GetProjectName(p.Ctx)
	if err != nil {
		return nil, err
	}
	pg, _ := p.Ctx.Values().GetBool("page")
	if pg {
		num, _ := p.Ctx.Values().GetInt(constant.PageNumQueryKey)
		size, _ := p.Ctx.Values().GetInt(constant.PageSizeQueryKey)
		return p.PlanService.Page(num, size, projectName, condition.TODO())
	} else {
		var page page.Page
		items, err := p.PlanService.List(projectName)
		if err != nil {
			return nil, err
		}
		page.Items = items
		page.Total = len(items)
		return &page, nil
	}
}

// Search Plan
// @Tags plans
// @Summary Search  plans
// @Description 过滤部署计划
// @Accept  json
// @Produce  json
// @Param conditions body condition.Conditions true "conditions"
// @Success 200 {object} page.Page
// @Security ApiKeyAuth
// @Router /plans/search [post]
func (p PlanController) PostSearch() (*page.Page, error) {
	pg, _ := p.Ctx.Values().GetBool("page")
	var conditions condition.Conditions
	if p.Ctx.GetContentLength() > 0 {
		if err := p.Ctx.ReadJSON(&conditions); err != nil {
			return nil, err
		}
	}
	projectName, err := sessionUtil.GetProjectName(p.Ctx)
	if err != nil {
		return nil, err
	}
	if pg {
		num, _ := p.Ctx.Values().GetInt(constant.PageNumQueryKey)
		size, _ := p.Ctx.Values().GetInt(constant.PageSizeQueryKey)
		return p.PlanService.Page(num, size, projectName, conditions)
	} else {
		var page page.Page
		items, err := p.PlanService.List(projectName)
		if err != nil {
			return nil, err
		}
		page.Items = items
		page.Total = len(items)
		return &page, nil
	}
}

// Get Plan
// @Tags plans
// @Summary Show a Plan
// @Description 获取单个部署计划
// @Accept  json
// @Produce  json
// @Param name path string true "部署计划名称"
// @Success 200 {object} dto.Plan
// @Security ApiKeyAuth
// @Router /plans/{name} [get]
func (p PlanController) GetBy(name string) (dto.Plan, error) {
	return p.PlanService.Get(name)
}

// Create Plan
// @Tags plans
// @Summary Create a plan
// @Description  创建部署计划
// @Accept  json
// @Produce  json
// @Param request body dto.PlanCreate true "request"
// @Success 200 {object} dto.Plan
// @Security ApiKeyAuth
// @Router /plans [post]
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

	operator := p.Ctx.Values().GetString("operator")
	go kolog.Save(operator, constant.CREATE_PLAN, req.Name)

	return p.PlanService.Create(req)
}

// Delete Plan
// @Tags plans
// @Summary Delete a plan
// @Description 删除部署计划
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Router /plans/{name} [delete]
func (p PlanController) DeleteBy(name string) error {
	operator := p.Ctx.Values().GetString("operator")
	go kolog.Save(operator, constant.DELETE_PLAN, name)

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

	operator := p.Ctx.Values().GetString("operator")
	delPlans := ""
	for _, item := range req.Items {
		delPlans += item.Name + ","
	}
	go kolog.Save(operator, constant.DELETE_PLAN, delPlans)

	return err
}

// Get Config
// @Tags plans
// @Summary Get vmConfigs
// @Description 获取虚拟机配置
// @Accept  json
// @Produce  json
// @Param name path string true "区域名称"
// @Success 200 {Array} []dto.PlanVmConfig
// @Security ApiKeyAuth
// @Router /configs/{regionName} [get]
func (p PlanController) GetConfigsBy(regionName string) ([]dto.PlanVmConfig, error) {
	return p.PlanService.GetConfigs(regionName)
}
