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
// @Description 获取项目成员列表
// @Accept  json
// @Produce  json
// @Param project path string true "项目名称"
// @Success 200 {object} page.Page
// @Security ApiKeyAuth
// @Router /project/{project}/members [get]
func (p ProjectMemberController) Get() (*page.Page, error) {
	projectName := p.Ctx.Params().GetString("project")
	pa, _ := p.Ctx.Values().GetBool("page")
	if pa {
		num, _ := p.Ctx.Values().GetInt(constant.PageNumQueryKey)
		size, _ := p.Ctx.Values().GetInt(constant.PageSizeQueryKey)
		return p.ProjectMemberService.Page(projectName, num, size)
	} else {
		var page page.Page
		return &page, nil
	}
}

func (p ProjectMemberController) GetBy(name string) (*dto.ProjectMember, error) {
	projectName := p.Ctx.Values().GetString("project")
	return p.ProjectMemberService.Get(name, projectName)
}

// Create ProjectMember
// @Tags projectMembers
// @Summary Create a projectMember
// @Description 授权成员到项目
// @Accept  json
// @Produce  json
// @Param request body dto.ProjectMemberCreate true "request"
// @Param project path string true "项目名称"
// @Success 200 {object} dto.ProjectMember
// @Security ApiKeyAuth
// @Router /project/{project}/members [post]
func (p ProjectMemberController) Post() ([]dto.ProjectMember, error) {
	projectName := p.Ctx.Params().GetString("project")
	var req dto.ProjectMemberCreate
	err := p.Ctx.ReadJSON(&req)
	if err != nil {
		return nil, err
	}

	result, err := p.ProjectMemberService.Create(projectName, req)
	if err != nil {
		return nil, err
	}

	operator := p.Ctx.Values().GetString("operator")
	go kolog.Save(operator, constant.BIND_PROJECT_MEMBER, projectName)

	return result, nil
}

// Delete Project Member
// @Tags projectMembers
// @Summary Delete Project Member
// @Description 取消项目人员授权
// @Accept  json
// @Produce  json
// @Param project path string true "项目名称"
// @Param name path string true "人员名称"
// @Security ApiKeyAuth
// @Router /project/{project}/members/{name} [delete]
func (p ProjectMemberController) DeleteBy(name string) error {
	projectName := p.Ctx.Params().GetString("project")
	return p.ProjectMemberService.Delete(name, projectName)
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
	//for _, item := range req.Items {
	//	delMembers += (item.Username + ",")
	//	delProject = item.ProjectName
	//}
	if req.Operation == "update" {
		go kolog.Save(operator, constant.UPDATE_PROJECT_MEMBER_ROLE, delProject+"-"+delMembers)
	} else {
		go kolog.Save(operator, constant.UNBIND_PROJECT_MEMBER, delProject+"-"+delMembers)
	}

	return err
}

func (p ProjectMemberController) GetUsers() (*dto.AddMemberResponse, error) {
	name := p.Ctx.URLParam("name")
	projectName := p.Ctx.Params().GetString("project")
	return p.ProjectMemberService.GetUsers(name, projectName)
}
