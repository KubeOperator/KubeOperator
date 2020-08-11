package controller

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/controller/page"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/service"
	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12/context"
)

type ProjectResourceController struct {
	Ctx                    context.Context
	ProjectResourceService service.ProjectResourceService
}

func NewProjectResourceController() *ProjectResourceController {
	return &ProjectResourceController{
		ProjectResourceService: service.NewProjectResourceService(),
	}
}

// List ProjectResource By ProjectName And ResourceType
// @Tags projectResources
// @Summary Show projectResources by projectName and resourceType
// @Description Show projectResources by projectName and resourceType
// @Accept  json
// @Produce  json
// @Form projectName string , resourceType string
// @Success 200 {object} page.Page
// @Security ApiKeyAuth
// @Router /project/resource/ [get]
func (p ProjectResourceController) Get() (page.Page, error) {
	pa, _ := p.Ctx.Values().GetBool("page")
	resourceType := p.Ctx.URLParam("resourceType")
	projectId := p.Ctx.URLParam("projectId")
	if pa {
		num, _ := p.Ctx.Values().GetInt(constant.PageNumQueryKey)
		size, _ := p.Ctx.Values().GetInt(constant.PageSizeQueryKey)
		return p.ProjectResourceService.PageByProjectIdAndType(num, size, projectId, resourceType)
	} else {
		var page page.Page
		return page, nil
	}
}

func (p ProjectResourceController) PostBatch() error {
	var req dto.ProjectResourceOp
	err := p.Ctx.ReadJSON(&req)
	if err != nil {
		return err
	}
	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		return err
	}
	err = p.ProjectResourceService.Batch(req)
	if err != nil {
		return err
	}
	return err
}

func (p ProjectResourceController) GetList() (interface{}, error) {
	resourceType := p.Ctx.URLParam("resourceType")
	projectName := p.Ctx.URLParam("projectName")
	return p.ProjectResourceService.GetResources(resourceType, projectName)
}
