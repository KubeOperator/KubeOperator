package controller

import (
	"errors"
	"fmt"

	"github.com/KubeOperator/KubeOperator/pkg/controller/condition"
	"github.com/KubeOperator/KubeOperator/pkg/util/nexus"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/controller/kolog"
	"github.com/KubeOperator/KubeOperator/pkg/controller/page"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/service"
	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12/context"
)

var (
	RegistryAlreadyExistsErr = errors.New("REGISTRY_ALREADY_EXISTS")
)

type SystemSettingController struct {
	Ctx                  context.Context
	SystemSettingService service.SystemSettingService
}

func NewSystemSettingController() *SystemSettingController {
	return &SystemSettingController{
		SystemSettingService: service.NewSystemSettingService(),
	}
}

// List SystemSettings
// @Tags SystemSetting
// @Summary Show all SystemSettings
// @Description 获取所有系统配置信息
// @Accept  json
// @Produce  json
// @Success 200 {object} dto.SystemSettingResult
// @Security ApiKeyAuth
// @Router /settings [get]
func (s SystemSettingController) Get() (interface{}, error) {
	item, err := s.SystemSettingService.List()
	if err != nil {
		return nil, err
	}
	return item, nil
}

// Get SystemSettings
// @Tags SystemSetting
// @Summary Show a SystemSettings
// @Description 获取单个应用配置的配置信息
// @Accept  json
// @Produce  json
// @Param name path string true "应用名称"
// @Success 200 {object} dto.SystemSettingResult
// @Security ApiKeyAuth
// @Router /settings/{name} [get]
func (s SystemSettingController) GetBy(name string) (interface{}, error) {
	item, err := s.SystemSettingService.ListByTab(name)
	if err != nil {
		return nil, err
	}
	return item, nil
}

// Create SystemSettings
// @Tags SystemSetting
// @Summary Create a SystemSetting
// @Description  创建一项配置
// @Accept  json
// @Produce  json
// @Param request body dto.SystemSettingCreate true "request"
// @Success 200 {object} []dto.SystemSetting
// @Security ApiKeyAuth
// @Router /settings [post]
func (s SystemSettingController) Post() ([]dto.SystemSetting, error) {
	var req dto.SystemSettingCreate
	err := s.Ctx.ReadJSON(&req)
	if err != nil {
		return nil, err
	}
	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		return nil, err
	}
	result, err := s.SystemSettingService.Create(req)
	if err != nil {
		return nil, err
	}

	operator := s.Ctx.Values().GetString("operator")
	go kolog.Save(operator, constant.CREATE_EMAIL, "-")

	return result, nil
}

// Check SystemSetting
// @Tags SystemSetting
// @Summary Check a SystemSetting
// @Description  检查配置是否可用
// @Accept  json
// @Produce  json
// @Param request body dto.SystemSettingCreate true "request"
// @Param name path string true "应用名称"
// @Success 200 {object} []dto.SystemSetting
// @Security ApiKeyAuth
// @Router /settings/check/{name} [post]
func (s SystemSettingController) PostCheckBy(typeName string) error {
	var req dto.SystemSettingCreate
	err := s.Ctx.ReadJSON(&req)
	if err != nil {
		return err
	}
	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		return err
	}
	err = s.SystemSettingService.CheckSettingByType(typeName, req)
	if err != nil {
		return err
	}
	return nil
}

// List Registry
// @Tags SystemSetting
// @Summary Show all Registry
// @Description 获取所有仓库信息
// @Accept  json
// @Produce  json
// @Success 200 {object} page.Page
// @Security ApiKeyAuth
// @Router /settings/registry [get]
func (s SystemSettingController) GetRegistry() (*page.Page, error) {
	p, _ := s.Ctx.Values().GetBool("page")
	if p {
		num, _ := s.Ctx.Values().GetInt(constant.PageNumQueryKey)
		size, _ := s.Ctx.Values().GetInt(constant.PageSizeQueryKey)
		return s.SystemSettingService.PageRegistry(num, size, condition.TODO())
	} else {
		var page page.Page
		items, err := s.SystemSettingService.ListRegistry(condition.TODO())
		if err != nil {
			return &page, err
		}
		page.Items = items
		page.Total = len(items)
		return &page, nil
	}

}

// Get Registry
// @Tags SystemSetting
// @Summary Show a Registry
// @Description 根据 ID 获取仓库信息
// @Accept  json
// @Produce  json
// @Param id path string true "ID"
// @Success 200 {object} dto.SystemRegistry
// @Security ApiKeyAuth
// @Router /settings/registry/{id} [get]
func (s SystemSettingController) GetRegistryBy(id string) (interface{}, error) {
	item, err := s.SystemSettingService.GetRegistryByID(id)
	if err != nil {
		return nil, err
	}
	return item, nil
}

// Create Registry
// @Tags SystemSetting
// @Summary Create a Registry
// @Description  创建仓库配置
// @Accept  json
// @Produce  json
// @Param request body dto.SystemSettingCreate true "request"
// @Success 200 {object} dto.SystemRegistry
// @Security ApiKeyAuth
// @Router /settings/registry [post]
func (s SystemSettingController) PostRegistry() (*dto.SystemRegistry, error) {
	var req dto.SystemRegistryCreate
	err := s.Ctx.ReadJSON(&req)
	if err != nil {
		return nil, err
	}
	item, _ := s.SystemSettingService.GetRegistryByArch(req.Architecture)
	if item.ID != "" {
		return nil, RegistryAlreadyExistsErr
	}

	operator := s.Ctx.Values().GetString("operator")
	go kolog.Save(operator, constant.CREATE_REGISTRY, req.Architecture)

	return s.SystemSettingService.CreateRegistry(req)
}

// Search Registry
// @Tags SystemSetting
// @Summary Search  Registry
// @Description 过滤仓库
// @Accept  json
// @Produce  json
// @Param conditions body condition.Conditions true "conditions"
// @Success 200 {object} page.Page
// @Security ApiKeyAuth
// @Router /settings/registry/search [post]
func (s SystemSettingController) PostRegistrySearch() (*page.Page, error) {
	var conditions condition.Conditions
	if s.Ctx.GetContentLength() > 0 {
		if err := s.Ctx.ReadJSON(&conditions); err != nil {
			return nil, err
		}
	}
	p, _ := s.Ctx.Values().GetBool("page")
	if p {
		num, _ := s.Ctx.Values().GetInt(constant.PageNumQueryKey)
		size, _ := s.Ctx.Values().GetInt(constant.PageSizeQueryKey)
		return s.SystemSettingService.PageRegistry(num, size, conditions)
	} else {
		var p page.Page
		items, err := s.SystemSettingService.ListRegistry(conditions)
		if err != nil {
			return &p, err
		}
		p.Items = items
		p.Total = len(items)
		return &p, nil
	}
}

// Update Registry
// @Tags SystemSetting
// @Summary Update a Registry
// @Description 更新仓库配置
// @Accept  json
// @Produce  json
// @Param request body dto.SystemRegistryUpdate true "request"
// @Param arch path string true "CPU 架构"
// @Success 200 {object} dto.SystemRegistry
// @Security ApiKeyAuth
// @Router /settings/registry/{arch} [patch]
func (s SystemSettingController) PatchRegistryBy(arch string) (*dto.SystemRegistry, error) {
	var req dto.SystemRegistryUpdate
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
	go kolog.Save(operator, constant.UPDATE_REGISTRY, req.Hostname)

	return s.SystemSettingService.UpdateRegistry(arch, req)
}

func (s SystemSettingController) PostRegistryBatch() error {
	var req dto.SystemRegistryBatchOp
	err := s.Ctx.ReadJSON(&req)
	if err != nil {
		return err
	}
	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		return err
	}
	err = s.SystemSettingService.BatchRegistry(req)
	if err != nil {
		return err
	}
	operator := s.Ctx.Values().GetString("operator")
	delCres := ""
	for _, item := range req.Items {
		delCres += (item.Architecture + ",")
	}
	go kolog.Save(operator, constant.DELETE_REGISTRY, delCres)
	return err
}

func (s SystemSettingController) PostRegistryCheckConn() error {
	var req dto.SystemRegistryConn
	err := s.Ctx.ReadJSON(&req)
	if err != nil {
		return err
	}
	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		return err
	}

	if err := nexus.CheckConn(req.Username, req.Password, fmt.Sprintf("%s://%s:%d", req.Protocol, req.Hostname, req.RepoPort)); err != nil {
		return err
	}
	return nil
}

// Delete Registry
// @Tags SystemSetting
// @Summary Delete a Registry
// @Description delete a  Registry by arch
// @Param arch path string true "CPU 架构"
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Router /settings/registry/{arch}/ [delete]
func (s SystemSettingController) DeleteRegistryBy(id string) error {
	go kolog.Save("Delete", constant.DELETE_REGISTRY, id)
	return s.SystemSettingService.DeleteRegistry(id)
}
