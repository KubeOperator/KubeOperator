package controller

import (
	"errors"
	"io/ioutil"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/controller/kolog"
	"github.com/KubeOperator/KubeOperator/pkg/controller/koregexp"
	"github.com/KubeOperator/KubeOperator/pkg/controller/page"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/service"
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
// @Param  pageNum  query  int  true "page number"
// @Param  pageSize  query  int  true "page size"
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
// @Param  pageNum  path  string  true "host name"
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
	if err := validate.RegisterValidation("koip", koregexp.CheckIpPattern); err != nil {
		return nil, err
	}
	if err := validate.Struct(req); err != nil {
		return nil, err
	}

	repos, err := h.SystemSettingService.GetLocalIPs()
	if err != nil {
		return nil, err
	}
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

	return &item, nil
}

func (h HostController) DeleteBy(name string) error {
	operator := h.Ctx.Values().GetString("operator")
	go kolog.Save(operator, constant.DELETE_HOST, name)

	return h.HostService.Delete(name)
}

func (h HostController) PostSync() error {
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

// Delete Hosts
// @Tags hosts
// @Summary delete host list
// @Description delete host list
// @Accept  json
// @Produce  json
// @Param request body dto.HostOp true "request"
// @Success 200 {object} dto.Host
// @Security ApiKeyAuth
// @Router /hosts/batch [post]
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
// @Accept  mpfd
// @Produce  json
// @Security ApiKeyAuth
// @Router /hosts/upload/ [post]
func (h HostController) PostUpload() error {
	f, _, err := h.Ctx.FormFile("file")
	if err != nil {
		return err
	}

	if sizeInterface, ok := f.(Size); ok {
		if sizeInterface.Size() > 10485760 {
			return errors.New("APP_HOST_IMPORT_FILE_SIZE_ERROR")
		}
	}

	bs, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	defer f.Close()

	operator := h.Ctx.Values().GetString("operator")
	go kolog.Save(operator, constant.UPLOAD_HOST, "-")

	return h.HostService.ImportHosts(bs)
}

type Size interface {
	Size() int64
}
