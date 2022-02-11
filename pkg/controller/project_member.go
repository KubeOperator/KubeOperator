package controller

import (
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

type ProjectMemberController struct {
	Ctx                  context.Context
	ProjectMemberService service.ProjectMemberService
}

func NewProjectMemberController() *ProjectMemberController {
	return &ProjectMemberController{
		ProjectMemberService: service.NewProjectMemberService(),
	}
}

// List ProjectMember By ProjectName
// @Tags projectMembers
// @Summary Show projectMembers by projectName
// @Description Show projectMembers by projectName
// @Accept  json
// @Produce  json
// @Param  pageNum  query  int  true "page number"
// @Param  pageSize  query  int  true "page size"
// @Param  project  header  string  true "project name"
// @Form projectName
// @Success 200 {object} page.Page
// @Security ApiKeyAuth
// @Router /project/members [get]
func (p ProjectMemberController) Get() (page.Page, error) {
	pa, _ := p.Ctx.Values().GetBool("page")

	projectNameUnDecode := p.Ctx.Request().Header.Get("project")
	projectName, err := url.QueryUnescape(projectNameUnDecode)
	if err != nil {
		var page page.Page
		return page, fmt.Errorf("decode error: %s", projectName)
	}

	if pa {
		num, _ := p.Ctx.Values().GetInt(constant.PageNumQueryKey)
		size, _ := p.Ctx.Values().GetInt(constant.PageSizeQueryKey)
		return p.ProjectMemberService.PageByProjectName(num, size, projectName)
	} else {
		var page page.Page
		return page, nil
	}
}

func (p ProjectMemberController) GetBy(name string) (*dto.ProjectMember, error) {
	projectNameUnDecode := p.Ctx.Request().Header.Get("project")
	projectName, err := url.QueryUnescape(projectNameUnDecode)
	if err != nil {
		return nil, fmt.Errorf("decode error: %s", projectName)
	}

	return p.ProjectMemberService.Get(name, projectName)
}

// Create ProjectMember
// @Tags projectMembers
// @Summary Create a projectMember
// @Description create a projectMember
// @Accept  json
// @Produce  json
// @Param request body dto.ProjectMemberCreate true "request"
// @Success 200 {object} dto.ProjectMember
// @Security ApiKeyAuth
// @Router /project/members/ [post]
func (p ProjectMemberController) Post() (*dto.ProjectMember, error) {
	var req dto.ProjectMemberCreate
	err := p.Ctx.ReadJSON(&req)
	if err != nil {
		return nil, err
	}

	result, err := p.ProjectMemberService.Create(req)
	if err != nil {
		return nil, err
	}

	operator := p.Ctx.Values().GetString("operator")
	go kolog.Save(operator, constant.BIND_PROJECT_MEMBER, req.ProjectName+"-"+req.Username)

	return result, nil
}

func (p ProjectMemberController) PostBatch() error {
	var req dto.ProjectMemberOP
	err := p.Ctx.ReadJSON(&req)
	if err != nil {
		return err
	}
	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		return err
	}
	err = p.ProjectMemberService.Batch(req)
	if err != nil {
		return err
	}

	operator := p.Ctx.Values().GetString("operator")
	delMembers, delProject := "", ""
	for _, item := range req.Items {
		delMembers += (item.Username + ",")
		delProject = item.ProjectName
	}
	if req.Operation == "update" {
		go kolog.Save(operator, constant.UPDATE_PROJECT_MEMBER_ROLE, delProject+"-"+delMembers)
	} else {
		go kolog.Save(operator, constant.UNBIND_PROJECT_MEMBER, delProject+"-"+delMembers)
	}

	return err
}

func (p ProjectMemberController) GetUsers() (dto.AddMemberResponse, error) {
	name := p.Ctx.URLParam("name")
	return p.ProjectMemberService.GetUsers(name)
}

func (p ProjectMemberController) GetRoles() ([]string, error) {
	var result []string
	result = append(result, constant.ProjectRoleProjectManager)
	result = append(result, constant.ProjectRoleClusterManager)
	return result, nil
}
