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

	sessionUser := p.Ctx.Values().Get("user")
	userId := getUserID(sessionUser)
	if userId == "UNRECOGNIZED_USER" {
		return page.Page{Items: []dto.Project{}, Total: 0}, errors.New("UNRECOGNIZED_USER")
	}
	hasPower, err := p.ProjectMemberService.CheckUserProjectPermissionByName(userId, []string{projectName})
	if !hasPower || err != nil {
		return page.Page{Items: []dto.Project{}, Total: 0}, errors.New("PERMISSION_DENIED")
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

	sessionUser := p.Ctx.Values().Get("user")
	userId := getUserID(sessionUser)
	if userId == "UNRECOGNIZED_USER" {
		return &dto.ProjectMember{}, errors.New("UNRECOGNIZED_USER")
	}
	hasPower, err := p.ProjectMemberService.CheckUserProjectPermissionByName(userId, []string{projectName})
	if !hasPower || err != nil {
		return &dto.ProjectMember{}, errors.New("PERMISSION_DENIED")
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
	if err := p.Ctx.ReadJSON(&req); err != nil {
		return nil, err
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return nil, err
	}

	result, err := p.ProjectMemberService.Create(req)
	if err != nil {
		return nil, err
	}

	sessionUser := p.Ctx.Values().Get("user")
	userId := getUserID(sessionUser)
	if userId == "UNRECOGNIZED_USER" {
		return &dto.ProjectMember{}, errors.New("UNRECOGNIZED_USER")
	}
	hasPower, err := p.ProjectMemberService.CheckUserProjectPermissionByName(userId, []string{req.ProjectName})
	if !hasPower || err != nil {
		return &dto.ProjectMember{}, errors.New("PERMISSION_DENIED")
	}

	operator := p.Ctx.Values().GetString("operator")
	go kolog.Save(operator, constant.BIND_PROJECT_MEMBER, req.ProjectName+"-"+req.Username)

	return result, nil
}

func (p ProjectMemberController) PostBatch() error {
	var req dto.ProjectMemberOP
	if err := p.Ctx.ReadJSON(&req); err != nil {
		return err
	}
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return err
	}
	if err := p.ProjectMemberService.Batch(req); err != nil {
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
	hasPower, err := p.ProjectMemberService.CheckUserProjectPermissionByName(userId, []string{req.Items[0].ProjectName})
	if !hasPower || err != nil {
		return errors.New("PERMISSION_DENIED")
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
