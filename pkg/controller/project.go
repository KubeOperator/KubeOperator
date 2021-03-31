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

type ProjectController struct {
	Ctx            context.Context
	ProjectService service.ProjectService
}

func NewProjectController() *ProjectController {
	return &ProjectController{
		ProjectService: service.NewProjectService(),
	}
}

// List Project
// @Tags projects
// @Summary Show all projects
// @Description Show projects
// @Accept  json
// @Produce  json
// @Success 200 {object} page.Page
// @Security ApiKeyAuth
// @Router /projects/ [get]
func (p ProjectController) Get() (page.Page, error) {
	pa, _ := p.Ctx.Values().GetBool("page")
	if pa {
		num, _ := p.Ctx.Values().GetInt(constant.PageNumQueryKey)
		size, _ := p.Ctx.Values().GetInt(constant.PageSizeQueryKey)
		sessionUser := p.Ctx.Values().Get("user")
		var userId string
		user, ok := sessionUser.(dto.SessionUser)
		if ok && !user.IsAdmin {
			userId = user.UserId
		} else {
			userId = ""
		}
		return p.ProjectService.Page(num, size, userId)
	} else {
		var page page.Page
		items, err := p.ProjectService.List()
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
// @Success 200 {object} dto.Project
// @Security ApiKeyAuth
// @Router /projects/{name}/ [get]
func (p ProjectController) GetBy(name string) (*dto.Project, error) {
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
	err = validate.Struct(req)
	if err != nil {
		return nil, err
	}
	result, err := p.ProjectService.Create(req)
	if err != nil {
		return result, err
	}

	operator := p.Ctx.Values().GetString("operator")
	go kolog.Save(operator, constant.CREATE_PROJECT, req.Name)

	return nil, err
}

// Update Project
// @Tags projects
// @Summary Update a project
// @Description Update a project
// @Accept  json
// @Produce  json
// @Param request body dto.ProjectUpdate true "request"
// @Success 200 {object} dto.Project
// @Security ApiKeyAuth
// @Router /projects/{name}/ [patch]
func (p ProjectController) PatchBy(name string) (*dto.Project, error) {
	var req dto.ProjectUpdate
	err := p.Ctx.ReadJSON(&req)
	if err != nil {
		return nil, err
	}
	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		return nil, err
	}
	result, err := p.ProjectService.Update(name, req)
	if err != nil {
		return nil, err
	}

	operator := p.Ctx.Values().GetString("operator")
	go kolog.Save(operator, constant.UPDATE_PROJECT_INFO, name)

	return result, nil
}

// Delete Project
// @Tags projects
// @Summary Delete a project
// @Description delete a  project by name
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Router /projects/{name}/ [delete]
func (p ProjectController) DeleteBy(name string) error {
	operator := p.Ctx.Values().GetString("operator")
	go kolog.Save(operator, constant.DELETE_PROJECT, name)

	return p.ProjectService.Delete(name)
}

func (p ProjectController) PostBatch() error {
	var req dto.ProjectOp
	err := p.Ctx.ReadJSON(&req)
	if err != nil {
		return err
	}
	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		return err
	}
	err = p.ProjectService.Batch(req)
	if err != nil {
		return err
	}

	operator := p.Ctx.Values().GetString("operator")
	delProjects := ""
	for _, item := range req.Items {
		delProjects += (item.Name + ",")
	}
	go kolog.Save(operator, constant.DELETE_PROJECT, delProjects)

	return err
}

func (p ProjectController) GetTree() ([]dto.ProjectResourceTree, error) {
	return p.ProjectService.GetResourceTree()
}
