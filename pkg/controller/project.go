package controller

import (
	"errors"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/controller/kolog"
	"github.com/KubeOperator/KubeOperator/pkg/controller/koregexp"
	"github.com/KubeOperator/KubeOperator/pkg/controller/page"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/service"
	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12/context"
)

type ProjectController struct {
	Ctx                  context.Context
	ProjectService       service.ProjectService
	ProjectMemberService service.ProjectMemberService
}

func NewProjectController() *ProjectController {
	return &ProjectController{
		ProjectService:       service.NewProjectService(),
		ProjectMemberService: service.NewProjectMemberService(),
	}
}

// List Project
// @Tags projects
// @Summary Show all projects
// @Description Show projects
// @Accept  json
// @Produce  json
// @Param  pageNum  query  int  true "page number"
// @Param  pageSize  query  int  true "page size"
// @Success 200 {object} page.Page
// @Security ApiKeyAuth
// @Router /projects/ [get]
func (p ProjectController) Get() (page.Page, error) {
	pa, _ := p.Ctx.Values().GetBool("page")
	sessionUser := p.Ctx.Values().Get("user")
	userId := getUserID(sessionUser)
	if userId == "UNRECOGNIZED_USER" {
		return page.Page{Items: []dto.Project{}, Total: 0}, errors.New("UNRECOGNIZED_USER")
	}

	if pa {
		num, _ := p.Ctx.Values().GetInt(constant.PageNumQueryKey)
		size, _ := p.Ctx.Values().GetInt(constant.PageSizeQueryKey)
		return p.ProjectService.Page(num, size, userId)
	} else {
		var page page.Page
		items, err := p.ProjectService.List(userId)
		if err != nil {
			return page, err
		}
		page.Items = items
		page.Total = len(items)
		return page, nil
	}
}

// Get Project
// @Tags projects
// @Summary Show a project
// @Description show a project by name
// @Accept  json
// @Produce  json
// @Param  name  path  int  true "page number"
// @Success 200 {object} dto.Project
// @Security ApiKeyAuth
// @Router /projects/{name}/ [get]
func (p ProjectController) GetBy(name string) (dto.Project, error) {
	sessionUser := p.Ctx.Values().Get("user")
	userId := getUserID(sessionUser)
	if userId == "UNRECOGNIZED_USER" {
		return dto.Project{}, errors.New("UNRECOGNIZED_USER")
	}
	hasPower, err := p.ProjectMemberService.CheckUserProjectPermissionByName(userId, []string{name})
	if !hasPower || err != nil {
		return dto.Project{}, errors.New("PERMISSION_DENIED")
	}

	return p.ProjectService.Get(name)
}

// Create Project
// @Tags projects
// @Summary Create a project
// @Description create a project
// @Accept  json
// @Produce  json
// @Param request body dto.ProjectCreate true "request"
// @Success 200 {object} dto.Project
// @Security ApiKeyAuth
// @Router /projects/ [post]
func (p ProjectController) Post() (*dto.Project, error) {
	var req dto.ProjectCreate
	err := p.Ctx.ReadJSON(&req)
	if err != nil {
		return nil, err
	}
	validate := validator.New()
	if err := validate.RegisterValidation("koname", koregexp.CheckNamePattern); err != nil {
		return nil, err
	}
	if err := validate.Struct(req); err != nil {
		return nil, err
	}
	result, err := p.ProjectService.Create(req)
	if err != nil {
		return result, err
	}

	go kolog.Save(p.Ctx, constant.CREATE_PROJECT, req.Name)

	return nil, err
}

// Update Project
// @Tags projects
// @Summary Update a project
// @Description Update a project
// @Accept  json
// @Produce  json
// @Param request body dto.ProjectUpdate true "request"
// @Param  name  path  string  true "project name"
// @Success 200 {object} dto.Project
// @Security ApiKeyAuth
// @Router /projects/{name}/ [patch]
func (p ProjectController) PatchBy(name string) (*dto.Project, error) {
	var req dto.ProjectUpdate
	if err := p.Ctx.ReadJSON(&req); err != nil {
		return nil, err
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return nil, err
	}

	sessionUser := p.Ctx.Values().Get("user")
	userId := getUserID(sessionUser)
	if userId == "UNRECOGNIZED_USER" {
		return &dto.Project{}, errors.New("UNRECOGNIZED_USER")
	}
	hasPower, err := p.ProjectMemberService.CheckUserProjectPermissionByName(userId, []string{name})
	if !hasPower || err != nil {
		return &dto.Project{}, errors.New("PERMISSION_DENIED")
	}

	result, err := p.ProjectService.Update(req)
	if err != nil {
		return nil, err
	}

	go kolog.Save(p.Ctx, constant.UPDATE_PROJECT_INFO, name)

	return &result, nil
}

func (p ProjectController) Delete(name string) error {
	sessionUser := p.Ctx.Values().Get("user")
	userId := getUserID(sessionUser)
	if userId == "UNRECOGNIZED_USER" {
		return errors.New("UNRECOGNIZED_USER")
	}

	go kolog.Save(p.Ctx, constant.DELETE_PROJECT, name)

	return p.ProjectService.Delete(name)
}

// Delete Projects
// @Tags projects
// @Summary Delete project list
// @Description delete  project list
// @Accept  json
// @Produce  json
// @Param request body dto.ProjectOp true "request"
// @Security ApiKeyAuth
// @Router /projects/batch [post]
func (p ProjectController) PostBatch() error {
	var req dto.ProjectOp
	if err := p.Ctx.ReadJSON(&req); err != nil {
		return err
	}
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return err
	}

	var delProjects string
	var projectNames []string
	for _, item := range req.Items {
		delProjects += (item.Name + ",")
		projectNames = append(projectNames, item.Name)
	}
	sessionUser := p.Ctx.Values().Get("user")
	userId := getUserID(sessionUser)
	if userId == "UNRECOGNIZED_USER" {
		return errors.New("UNRECOGNIZED_USER")
	}
	hasPower, err := p.ProjectMemberService.CheckUserProjectPermissionByName(userId, projectNames)
	if !hasPower || err != nil {
		return errors.New("PERMISSION_DENIED")
	}

	if err := p.ProjectService.Batch(req); err != nil {
		return err
	}

	go kolog.Save(p.Ctx, constant.DELETE_PROJECT, delProjects)

	return err
}

func getUserID(sessionUser interface{}) string {
	user, ok := sessionUser.(dto.SessionUser)
	if !ok {
		return "UNRECOGNIZED_USER"
	}
	if user.IsAdmin {
		return ""
	}
	return user.UserId
}
