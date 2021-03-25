package controller

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/controller/condition"
	"github.com/KubeOperator/KubeOperator/pkg/controller/kolog"
	"github.com/KubeOperator/KubeOperator/pkg/controller/page"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/service"
	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12/context"
)

type IpController struct {
	Ctx       context.Context
	IpService service.IpService
}

func NewIpController() *IpController {
	return &IpController{
		IpService: service.NewIpService(),
	}
}

// List IP
// @Tags ips
// @Summary Show ips by ipPoolName
// @Description Show ips by ipPoolName
// @Accept  json
// @Produce  json
// @Success 200 {object} page.Page
// @Security ApiKeyAuth
// @Router /ippools/{name}/ips/ [get]
func (i IpController) Get() (*page.Page, error) {
	p, _ := i.Ctx.Values().GetBool("page")
	ipPoolName := i.Ctx.Params().GetString("name")
	if p {
		num, _ := i.Ctx.Values().GetInt(constant.PageNumQueryKey)
		size, _ := i.Ctx.Values().GetInt(constant.PageSizeQueryKey)
		return i.IpService.Page(num, size, ipPoolName, condition.TODO())
	} else {
		var p page.Page
		items, err := i.IpService.List(ipPoolName, condition.TODO())
		if err != nil {
			return &p, err
		}
		p.Items = items
		p.Total = len(items)
		return &p, nil
	}
}

func (i IpController) PostSearch() (*page.Page, error) {
	ipPoolName := i.Ctx.Params().GetString("name")
	var conditions condition.Conditions
	if i.Ctx.GetContentLength() > 0 {
		if err := i.Ctx.ReadJSON(&conditions); err != nil {
			return nil, err
		}
	}
	p, _ := i.Ctx.Values().GetBool("page")
	if p {
		num, _ := i.Ctx.Values().GetInt(constant.PageNumQueryKey)
		size, _ := i.Ctx.Values().GetInt(constant.PageSizeQueryKey)
		return i.IpService.Page(num, size, ipPoolName, conditions)
	} else {
		var p page.Page
		items, err := i.IpService.List(ipPoolName, condition.TODO())
		if err != nil {
			return &p, err
		}
		p.Items = items
		p.Total = len(items)
		return &p, nil
	}
}

// Create Ip
// @Tags ips
// @Summary Create a Ip
// @Description create a Ip
// @Accept  json
// @Produce  json
// @Param request body dto.IpCreate true "request"
// @Success 200 {object} dto.Ip
// @Security ApiKeyAuth
// @Router /ippools/{name}/ips [post]
func (i IpController) Post() error {
	var req dto.IpCreate
	err := i.Ctx.ReadJSON(&req)
	if err != nil {
		return err
	}
	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		return err
	}
	return i.IpService.Create(req, nil)
}

func (i IpController) PostBatch() error {
	var req dto.IpOp
	err := i.Ctx.ReadJSON(&req)
	if err != nil {
		return err
	}
	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		return err
	}
	return i.IpService.Batch(req)
}

// Update Ip
// @Tags ips
// @Summary Update a Ip
// @Description Update a Ip
// @Accept  json
// @Produce  json
// @Param request body dto.IpUpdate true "request"
// @Success 200 {object} dto.Ip
// @Security ApiKeyAuth
// @Router /ippools/{name}/ips  [patch]
func (i IpController) Patch() (*dto.Ip, error) {
	var req dto.IpUpdate
	err := i.Ctx.ReadJSON(&req)
	if err != nil {
		return nil, err
	}
	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		return nil, err
	}
	return i.IpService.Update(req)
}

func (i IpController) PatchSync() error {
	ipPoolName := i.Ctx.Params().Get("name")
	return i.IpService.Sync(ipPoolName)
}

// Delete Ip
// @Tags ips
// @Summary Delete a Ip
// @Description delete a ip by address
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Router /ippools/{name}/ips/{address} [delete]
func (i IpController) DeleteBy(address string) error {
	operator := i.Ctx.Values().GetString("operator")
	go kolog.Save(operator, constant.DELETE_IP, address)
	return i.IpService.Delete(address)
}
