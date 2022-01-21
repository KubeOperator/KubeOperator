package controller

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
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

func (i IpController) Get() (page.Page, error) {
	p, _ := i.Ctx.Values().GetBool("page")
	ipPoolName := i.Ctx.Params().GetString("name")
	if p {
		num, _ := i.Ctx.Values().GetInt(constant.PageNumQueryKey)
		size, _ := i.Ctx.Values().GetInt(constant.PageSizeQueryKey)
		return i.IpService.Page(num, size, ipPoolName)
	} else {
		return page.Page{}, nil
	}
}

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

func (i IpController) PostSync() error {
	var req dto.IpSync
	err := i.Ctx.ReadJSON(&req)
	if err != nil {
		return err
	}
	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		return err
	}
	return i.IpService.Sync(req.IpPoolName)
}
