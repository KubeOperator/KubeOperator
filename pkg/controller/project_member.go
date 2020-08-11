package controller

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/controller/page"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/permission"
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
// @Form projectName
// @Success 200 {object} page.Page
// @Security ApiKeyAuth
// @Router /project/members/ [get]
func (p ProjectMemberController) Get() (page.Page, error) {
	pa, _ := p.Ctx.Values().GetBool("page")
	projectName := p.Ctx.URLParam("projectName")
	if pa {
		num, _ := p.Ctx.Values().GetInt(constant.PageNumQueryKey)
		size, _ := p.Ctx.Values().GetInt(constant.PageSizeQueryKey)
		return p.ProjectMemberService.PageByProjectName(num, size, projectName)
	} else {
		var page page.Page
		return page, nil
	}
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
	return err
}

func (p ProjectMemberController) GetUsers() (dto.AddMemberResponse, error) {
	name := p.Ctx.URLParam("name")
	return p.ProjectMemberService.GetUsers(name)
}

func (p ProjectMemberController) GetRoles() ([]string, error) {
	var result []string
	result = append(result, permission.CLUSTERMANAGER)
	result = append(result, permission.PROJECTMANAGER)
	return result, nil
}
