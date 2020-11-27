package controller

import (
	"errors"
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/controller/log_save"
	"github.com/KubeOperator/KubeOperator/pkg/controller/page"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/service"
	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12/context"
	"io/ioutil"
)

var (
	HostAlreadyExistsErr = "HOST_ALREADY_EXISTS"
	SystemIpNotFound     = "SYSTEM_IP_NOT_FOUND"
)

type HostController struct {
	Ctx                  context.Context
	HostService          service.HostService
	SystemSettingService service.SystemSettingService
}

func NewHostController() *HostController {
	return &HostController{
		HostService:          service.NewHostService(),
		SystemSettingService: service.NewSystemSettingService(),
	}
}

// List Host
// @Tags hosts
// @Summary Show all hosts
// @Description Show hosts
// @Accept  json
// @Produce  json
// @Success 200 {object} page.Page
// @Security ApiKeyAuth
// @Router /hosts/ [get]
func (h HostController) Get() (page.Page, error) {

	p, _ := h.Ctx.Values().GetBool("page")
	if p {
		num, _ := h.Ctx.Values().GetInt(constant.PageNumQueryKey)
		size, _ := h.Ctx.Values().GetInt(constant.PageSizeQueryKey)
		return h.HostService.Page(num, size)
	} else {
		var page page.Page
		projectName := h.Ctx.URLParam("projectName")
		items, err := h.HostService.List(projectName)
		if err != nil {
			return page, err
		}
		page.Items = items
		page.Total = len(items)
		return page, nil
	}
}

// Get Host
// @Tags hosts
// @Summary Show a host
// @Description show a host by name
// @Accept  json
// @Produce  json
// @Success 200 {object} dto.Host
// @Security ApiKeyAuth
// @Router /hosts/{name}/ [get]
func (h HostController) GetBy(name string) (*dto.Host, error) {
	ho, err := h.HostService.Get(name)
	if err != nil {
		return nil, err
	}
	return &ho, nil

}

// Create Host
// @Tags hosts
// @Summary Create a host
// @Description create a host
// @Accept  json
// @Produce  json
// @Param request body dto.HostCreate true "request"
// @Success 200 {object} dto.Host
// @Security ApiKeyAuth
// @Router /hosts/ [post]
func (h HostController) Post() (*dto.Host, error) {
	var req dto.HostCreate
	err := h.Ctx.ReadJSON(&req)
	if err != nil {
		return nil, err
	}
	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		return nil, err
	}

	localIp, err := h.SystemSettingService.Get("ip")
	if err != nil {
		return nil, errors.New(SystemIpNotFound)
	}
	if localIp.Value == req.Ip {
		return nil, errors.New("IS_LOCAL_HOST")
	}
	item, _ := h.HostService.Get(req.Name)
	if item.ID != "" {
		return nil, errors.New(HostAlreadyExistsErr)
	}
	item, err = h.HostService.Create(req)
	if err != nil {
		return nil, err
	}

	operator := h.Ctx.Values().GetString("operator")
	go log_save.LogSave(operator, constant.CREATE_HOST, req.Name)

	return &item, nil
}

// Delete Host
// @Tags hosts
// @Summary Delete a host
// @Description delete a host by name
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Router /hosts/{name}/ [delete]
func (h HostController) Delete(name string) error {
	operator := h.Ctx.Values().GetString("operator")
	go log_save.LogSave(operator, constant.DELETE_HOST, name)

	return h.HostService.Delete(name)
}

func (h HostController) PostSyncBy(name string) (dto.Host, error) {
	return h.HostService.Sync(name)
}

func (h HostController) PostBatch() error {
	var req dto.HostOp
	err := h.Ctx.ReadJSON(&req)
	if err != nil {
		return err
	}
	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		return err
	}
	err = h.HostService.Batch(req)
	if err != nil {
		return err
	}

	operator := h.Ctx.Values().GetString("operator")
	delHost := ""
	for _, item := range req.Items {
		delHost += (item.Name + ",")
	}
	go log_save.LogSave(operator, constant.DELETE_HOST, delHost)

	return err
}

// Download Host Template File
// @Tags hosts
// @Summary Download Host Template File
// @Description download template file for import hosts
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Router /hosts/template/ [get]
func (h HostController) GetTemplate() error {
	err := h.HostService.DownloadTemplateFile()
	if err != nil {
		return err
	}
	err = h.Ctx.SendFile("demo.xlsx", "./demo.xlsx")
	if err != nil {
		return err
	}
	return nil
}

// Upload File for import
// @Tags hosts
// @Summary Upload File for import
// @Description Upload File for import hosts
// @Accept  xlsx
// @Produce  json
// @Security ApiKeyAuth
// @Router /hosts/upload/ [post]
func (h HostController) PostUpload() (*dto.ImportHostResponse, error) {
	f, _, err := h.Ctx.FormFile("file")
	if err != nil {
		return nil, err
	}
	bs, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return h.HostService.ImportHosts(bs)
}
