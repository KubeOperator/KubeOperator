package controller

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/controller/kolog"
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
// @Param  pageNum  query  int  true "page number"
// @Param  pageSize  query  int  true "page size"
// @Param  resourceType  query  string  true "resourceType enums"  Enums(HOST, BACKUP_ACCOUNT)
// @Param  project  header  string  true "project name"
// @Success 200 {object} page.Page
// @Security ApiKeyAuth
// @Router /project/resources [get]
func (p ProjectResourceController) Get() (*page.Page, error) {
	pa, _ := p.Ctx.Values().GetBool("page")
	resourceType := p.Ctx.URLParam("resourceType")
	projectName := p.Ctx.Values().GetString("project")
	if pa {
		num, _ := p.Ctx.Values().GetInt(constant.PageNumQueryKey)
		size, _ := p.Ctx.Values().GetInt(constant.PageSizeQueryKey)
		return p.ProjectResourceService.Page(num, size, projectName, resourceType)
	} else {
		return nil, nil
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

	operator := p.Ctx.Values().GetString("operator")
	go saveResourceBindLogs(operator, req)

	return err
}

func (p ProjectResourceController) GetList() (interface{}, error) {
	resourceType := p.Ctx.URLParam("resourceType")
	projectName := p.Ctx.Values().GetString("project")
	return p.ProjectResourceService.GetResources(resourceType, projectName)
}

func saveResourceBindLogs(operator string, req dto.ProjectResourceOp) {
	resources := ""
	typeStr := ""
	for _, item := range req.Items {
		typeStr = item.ResourceType
		resources += (item.ResourceName + ",")
	}
	if req.Operation == "create" {
		switch typeStr {
		case "PLAN":
			go kolog.Save(operator, constant.BIND_PROJECT_RESOURCE_PLAN, resources)
		case "BACKUP_ACCOUNT":
			go kolog.Save(operator, constant.BIND_PROJECT_RESOURCE_BACKUP, resources)
		case "HOST":
			go kolog.Save(operator, constant.BIND_PROJECT_RESOURCE_HOST, resources)
		}
	} else {
		switch typeStr {
		case "PLAN":
			go kolog.Save(operator, constant.UNBIND_PROJECT_RESOURCE_PLAN, resources)
		case "BACKUP_ACCOUNT":
			go kolog.Save(operator, constant.UNBIND_PROJECT_RESOURCE_BACKUP, resources)
		case "HOST":
			go kolog.Save(operator, constant.UNBIND_PROJECT_RESOURCE_HOST, resources)
		}
	}
}
