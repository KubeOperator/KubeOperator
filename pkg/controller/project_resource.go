package controller

import (
	"errors"
	"fmt"
	"net/url"

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
	ProjectMemberService   service.ProjectMemberService
}

func NewProjectResourceController() *ProjectResourceController {
	return &ProjectResourceController{
		ProjectResourceService: service.NewProjectResourceService(),
		ProjectMemberService:   service.NewProjectMemberService(),
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

	projectNameUnDecode := p.Ctx.Request().Header.Get("project")
	projectName, err := url.QueryUnescape(projectNameUnDecode)
	if err != nil {
		return nil, fmt.Errorf("decode error: %s", projectName)
	}

	sessionUser := p.Ctx.Values().Get("user")
	userId := getUserID(sessionUser)
	if userId == "UNRECOGNIZED_USER" {
		return &page.Page{Items: []dto.Project{}, Total: 0}, errors.New("UNRECOGNIZED_USER")
	}
	hasPower, err := p.ProjectMemberService.CheckUserProjectPermissionByName(userId, []string{projectName})
	if !hasPower || err != nil {
		return &page.Page{Items: []dto.Project{}, Total: 0}, errors.New("PERMISSION_DENIED")
	}

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
	if err := p.Ctx.ReadJSON(&req); err != nil {
		return err
	}
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return err
	}
	if err := p.ProjectResourceService.Batch(req); err != nil {
		return err
	}

	sessionUser := p.Ctx.Values().Get("user")
	userId := getUserID(sessionUser)
	if userId == "UNRECOGNIZED_USER" {
		return errors.New("UNRECOGNIZED_USER")
	}
	if len(req.Items) == 0 {
		return nil
	}
	hasPower, err := p.ProjectMemberService.CheckUserProjectPermissionByID(userId, []string{req.Items[0].ProjectID})
	if !hasPower || err != nil {
		return errors.New("PERMISSION_DENIED")
	}

	operator := p.Ctx.Values().GetString("operator")
	go saveResourceBindLogs(operator, req)

	return err
}

func (p ProjectResourceController) GetList() (interface{}, error) {
	resourceType := p.Ctx.URLParam("resourceType")
	projectNameUnDecode := p.Ctx.Request().Header.Get("project")
	projectName, err := url.QueryUnescape(projectNameUnDecode)
	if err != nil {
		return nil, fmt.Errorf("decode error: %s", projectName)
	}

	sessionUser := p.Ctx.Values().Get("user")
	userId := getUserID(sessionUser)
	if userId == "UNRECOGNIZED_USER" {
		return nil, errors.New("UNRECOGNIZED_USER")
	}
	hasPower, err := p.ProjectMemberService.CheckUserProjectPermissionByName(userId, []string{projectName})
	if !hasPower || err != nil {
		return nil, errors.New("PERMISSION_DENIED")
	}

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
