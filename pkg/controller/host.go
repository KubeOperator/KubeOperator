package controller

import (
	"errors"
	"github.com/KubeOperator/KubeOperator/pkg/controller/condition"
	"io/ioutil"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/controller/kolog"
	"github.com/KubeOperator/KubeOperator/pkg/controller/page"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/service"
	sessionUtil "github.com/KubeOperator/KubeOperator/pkg/util/session"
	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12/context"
)

var (
	HostAlreadyExistsErr     = "HOST_ALREADY_EXISTS"
	SystemRegistryIpNotFound = "SYSTEM_REGISTRY_IP_NOT_FOUND"
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
func (h *HostController) Get() (*page.Page, error) {
	p, _ := h.Ctx.Values().GetBool("page")
	profile, err := sessionUtil.GetUser(h.Ctx)
	if err != nil {
		return nil, err
	}
	if p {
		num, _ := h.Ctx.Values().GetInt(constant.PageNumQueryKey)
		size, _ := h.Ctx.Values().GetInt(constant.PageSizeQueryKey)
		return h.HostService.Page(num, size, profile.User.CurrentProject, condition.TODO())
	} else {
		var p page.Page
		items, err := h.HostService.List(profile.User.CurrentProject, condition.TODO())
		if err != nil {
			return &p, err
		}
		p.Items = items
		p.Total = len(items)
		return &p, nil
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
func (h *HostController) GetBy(name string) (*dto.Host, error) {
	return h.HostService.Get(name)
}

func (h *HostController) PostSearch() (*page.Page, error) {
	var conditions condition.Conditions
	profile, err := sessionUtil.GetUser(h.Ctx)
	if err != nil {
		return nil, err
	}
	if h.Ctx.GetContentLength() > 0 {
		if err := h.Ctx.ReadJSON(&conditions); err != nil {
			return nil, err
		}
	}
	p, _ := h.Ctx.Values().GetBool("page")
	if p {
		num, _ := h.Ctx.Values().GetInt(constant.PageNumQueryKey)
		size, _ := h.Ctx.Values().GetInt(constant.PageSizeQueryKey)
		return h.HostService.Page(num, size, profile.User.CurrentProject, conditions)
	} else {
		var p page.Page
		items, err := h.HostService.List(profile.User.CurrentProject, conditions)
		if err != nil {
			return &p, err
		}
		p.Items = items
		p.Total = len(items)
		return &p, nil
	}
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
func (h *HostController) Post() (*dto.Host, error) {
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

	repos, err := h.SystemSettingService.GetLocalIPs()
	isExit := false
	for _, repo := range repos {
		if repo.Hostname == req.Ip {
			isExit = true
		}
	}
	if isExit {
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
	go kolog.Save(operator, constant.CREATE_HOST, req.Name)

	return item, nil
}

// Delete Host
// @Tags hosts
// @Summary Delete a host
// @Description delete a host by name
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Router /hosts/{name}/ [delete]
func (h *HostController) DeleteBy(name string) error {
	operator := h.Ctx.Values().GetString("operator")
	go kolog.Save(operator, constant.DELETE_HOST, name)

	return h.HostService.Delete(name)
}

func (h *HostController) PostSync() error {
	var req []dto.HostSync
	err := h.Ctx.ReadJSON(&req)
	if err != nil {
		return err
	}

	var hostStr string
	for _, host := range req {
		hostStr += (host.HostName + ",")
	}
	operator := h.Ctx.Values().GetString("operator")
	go kolog.Save(operator, constant.SYNC_HOST_LIST, hostStr)

	return h.HostService.SyncList(req)
}

func (h *HostController) PostBatch() error {
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
	go kolog.Save(operator, constant.DELETE_HOST, delHost)

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
func (h *HostController) GetTemplate() error {
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
// @Accept  mpfd
// @Produce  json
// @Security ApiKeyAuth
// @Router /hosts/upload/ [post]
func (h *HostController) PostUpload() error {
	f, _, err := h.Ctx.FormFile("file")
	if err != nil {
		return err
	}
	bs, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	defer f.Close()
	return h.HostService.ImportHosts(bs)
}
