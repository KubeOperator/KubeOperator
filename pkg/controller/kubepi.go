package controller

import (
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/service"
	"github.com/kataras/iris/v12/context"
)

type KubePiController struct {
	Ctx           context.Context
	KubePiService service.KubepiService
}

func NewKubePiController() *KubePiController {
	return &KubePiController{
		KubePiService: service.NewKubepiService(),
	}
}

func (u *KubePiController) GetUser() (interface{}, error) {
	users, err := u.KubePiService.GetKubePiUser()
	return users, err
}

func (p KubePiController) PostBind() error {
	var req dto.BindKubePI
	err := p.Ctx.ReadJSON(&req)
	if err != nil {
		return err
	}

	if err := p.KubePiService.BindKubePi(req); err != nil {
		return err
	}

	// operator := p.Ctx.Values().GetString("operator")
	// go kolog.Save(operator, constant.BIND_PROJECT_MEMBER, projectName)
	return nil
}

func (p KubePiController) PostSearch() (*dto.BindResponse, error) {
	var req dto.SearchBind
	err := p.Ctx.ReadJSON(&req)
	if err != nil {
		return nil, err
	}

	bind, err := p.KubePiService.GetKubePiBind(req)
	if err != nil {
		return nil, err
	}
	return bind, nil
}

func (p KubePiController) PostCheckConn() error {
	var req dto.CheckConn
	err := p.Ctx.ReadJSON(&req)
	if err != nil {
		return err
	}

	return p.KubePiService.CheckConn(req)
}
