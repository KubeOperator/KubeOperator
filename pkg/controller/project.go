package controller

import (
	"github.com/KubeOperator/KubeOperator/pkg/auth"
	"github.com/KubeOperator/KubeOperator/pkg/constant"
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

func (p ProjectController) Get() (page.Page, error) {

	pa, _ := p.Ctx.Values().GetBool("page")
	if pa {
		num, _ := p.Ctx.Values().GetInt(constant.PageNumQueryKey)
		size, _ := p.Ctx.Values().GetInt(constant.PageSizeQueryKey)
		sessionUser := p.Ctx.Values().Get("user")
		var userId string
		user, ok := sessionUser.(auth.SessionUser)
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

func (p ProjectController) GetBy(name string) (dto.Project, error) {
	return p.ProjectService.Get(name)
}

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
	return nil, err
}

func (p ProjectController) PatchBy(name string) (dto.Project, error) {
	var req dto.ProjectUpdate
	err := p.Ctx.ReadJSON(&req)
	if err != nil {
		return dto.Project{}, err
	}
	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		return dto.Project{}, err
	}
	return p.ProjectService.Update(req)
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
	return err
}
