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

type NtpServerController struct {
	Ctx              context.Context
	NtpServerService service.NtpServerService
}

func NewNtpServerController() *NtpServerController {
	return &NtpServerController{
		NtpServerService: service.NewNtpServerService(),
	}
}

// List NtpServer
// @Tags ntpServer
// @Summary Show ntpServer
// @Description 获取用户列表
// @Accept  json
// @Produce  json
// @Success 200 {object} page.Page
// @Security ApiKeyAuth
// @Router /ntpServer [get]
func (u *NtpServerController) Get() (*page.Page, error) {
	p, _ := u.Ctx.Values().GetBool("page")
	if p {
		num, _ := u.Ctx.Values().GetInt(constant.PageNumQueryKey)
		size, _ := u.Ctx.Values().GetInt(constant.PageSizeQueryKey)
		return u.NtpServerService.Page(num, size)
	} else {
		var p page.Page
		items, err := u.NtpServerService.List()
		if err != nil {
			return &p, err
		}
		p.Items = items
		p.Total = len(items)
		return &p, nil
	}
}

// Create NtpServers
// @Tags NtpServer
// @Summary Create a NtpServer
// @Description  创建ntpServer
// @Accept  json
// @Produce  json
// @Param request body dto.NtpServerCreate true "request"
// @Success 200 {object} dto.NtpServer
// @Security ApiKeyAuth
// @Router /ntp [post]
func (s NtpServerController) Post() (*dto.NtpServer, error) {
	var req dto.NtpServerCreate
	err := s.Ctx.ReadJSON(&req)
	if err != nil {
		return nil, err
	}
	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		return nil, err
	}
	result, err := s.NtpServerService.Create(req)
	if err != nil {
		return nil, err
	}

	operator := s.Ctx.Values().GetString("operator")
	go kolog.Save(operator, constant.CREATE_NTP, req.Name)

	return result, nil
}

// Update Registry
// @Tags NtpServer
// @Summary Update a Registry
// @Description 更新NtpServer
// @Accept  json
// @Produce  json
// @Param request body dto.NtpServerUpdate true "request"
// @Param name path string true "名称"
// @Success 200 {object} dto.NtpServer
// @Security ApiKeyAuth
// @Router /ntp/{name} [patch]
func (s NtpServerController) PatchBy(name string) (*dto.NtpServer, error) {
	var req dto.NtpServerUpdate
	err := s.Ctx.ReadJSON(&req)
	if err != nil {
		return nil, err
	}
	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		return nil, err
	}
	operator := s.Ctx.Values().GetString("operator")
	go kolog.Save(operator, constant.UPDATE_NTP, name)

	return s.NtpServerService.Update(name, req)
}

// Delete Registry
// @Tags NtpServer
// @Summary Delete a Registry
// @Description delete a  Registry by name
// @Param name path string true "CPU 架构"
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Router /ntp/{name}/ [delete]
func (s NtpServerController) DeleteBy(name string) error {
	go kolog.Save("Delete", constant.DELETE_NTP, name)
	return s.NtpServerService.Delete(name)
}
