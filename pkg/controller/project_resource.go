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

// List ProjectResource By ResourceType
// @Tags projectResources
// @Summary Show projectResources by resourceType
// @Description 分页获取项目资源列表
// @Accept  json
// @Produce  json
// @Param resourceType query string true  "资源类型（HOST,PLAN,BACKUP_ACCOUNT）"
// @Param project path string true "项目名称"
// @Success 200 {object} page.Page
// @Security ApiKeyAuth
// @Router /projects/{project}/resources [get]
func (p ProjectResourceController) Get() (*page.Page, error) {
	pa, _ := p.Ctx.Values().GetBool("page")
	resourceType := p.Ctx.URLParam("resourceType")
	projectName := p.Ctx.Params().GetString("project")
	if pa {
		num, _ := p.Ctx.Values().GetInt(constant.PageNumQueryKey)
		size, _ := p.Ctx.Values().GetInt(constant.PageSizeQueryKey)
		return p.ProjectResourceService.Page(num, size, projectName, resourceType)
	} else {
		return p.ProjectResourceService.List(projectName, resourceType)
	}
}

// Create ProjectResource
// @Tags projectResources
// @Summary Create a projectResource
// @Description 授权资源到项目
// @Accept  json
// @Produce  json
// @Param request body dto.ProjectResourceCreate true "request"
// @Param project path string true "项目名称"
// @Success 200 {object} dto.ProjectResource
// @Security ApiKeyAuth
// @Router /projects/{project}/resources [post]
func (p ProjectResourceController) Post() ([]dto.ProjectResource, error) {
	projectName := p.Ctx.Params().GetString("project")

	var req dto.ProjectResourceCreate
	err := p.Ctx.ReadJSON(&req)
	if err != nil {
		return nil, err
	}
	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		return nil, err
	}
	return p.ProjectResourceService.Create(projectName, req)
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

	operator := p.Ctx.Params().GetString("operator")
	go saveResourceBindLogs(operator, req)

	return err
}

// Get Project Resources
// @Tags projectResources
// @Summary Get projectResources
// @Description 获取能添加到项目的资源
// @Accept  json
// @Produce  json
// @Param project path string true "项目名称"
// @Success 200 {object} interface{}
// @Security ApiKeyAuth
// @Router /projects/{project}/resources/list [get]
func (p ProjectResourceController) GetList() (interface{}, error) {
	resourceType := p.Ctx.URLParam("resourceType")
	projectName := p.Ctx.Params().GetString("project")
	return p.ProjectResourceService.GetResources(resourceType, projectName)
}

// Delete Project Resource
// @Tags projectResources
// @Summary Delete projectResource
// @Description 取消项目资源授权
// @Accept  json
// @Produce  json
// @Param project path string true "项目名称"
// @Param name path string true "资源名称"
// @Security ApiKeyAuth
// @Router /projects/{project}/resources/{name} [delete]
func (p ProjectResourceController) DeleteBy(name string) error {
	resourceType := p.Ctx.URLParam("resourceType")
	projectName := p.Ctx.Params().GetString("project")
	return p.ProjectResourceService.Delete(name, resourceType, projectName)
}

func saveResourceBindLogs(operator string, req dto.ProjectResourceOp) {
	resources := ""
	typeStr := ""
	for _, item := range req.Items {
		typeStr = item.ResourceType
		resources += item.ResourceName + ","
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
