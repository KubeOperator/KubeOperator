package controller

import (
	"errors"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/controller/koregexp"
	"github.com/KubeOperator/KubeOperator/pkg/controller/page"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/service"
	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12/context"
)

type IpPoolController struct {
	Ctx           context.Context
	IpPoolService service.IpPoolService
}

func NewIpPoolController() *IpPoolController {
	return &IpPoolController{
		IpPoolService: service.NewIpPoolService(),
	}
}

func (i IpPoolController) Get() (page.Page, error) {
	p, _ := i.Ctx.Values().GetBool("page")
	if p {
		num, _ := i.Ctx.Values().GetInt(constant.PageNumQueryKey)
		size, _ := i.Ctx.Values().GetInt(constant.PageSizeQueryKey)
		return i.IpPoolService.Page(num, size)
	} else {
		var p page.Page
		items, err := i.IpPoolService.List()
		if err != nil {
			return page.Page{}, err
		}
		p.Items = items
		p.Total = len(items)
		return p, nil
	}
}

func (i IpPoolController) Post() (*dto.IpPool, error) {
	var req dto.IpPoolCreate
	err := i.Ctx.ReadJSON(&req)
	if err != nil {
		return nil, err
	}
	validate := validator.New()
	if err := validate.RegisterValidation("koip", koregexp.CheckIpPattern); err != nil {
		return nil, err
	}
	if err := validate.RegisterValidation("koname", koregexp.CheckNamePattern); err != nil {
		return nil, err
	}
	if err := validate.Struct(req); err != nil {
		return nil, err
	}
	item, _ := i.IpPoolService.Get(req.Name)
	if item.ID != "" {
		return nil, errors.New("NAME_EXISTS")
	}
	item, err = i.IpPoolService.Create(req)
	if err != nil {
		return nil, err
	}
	return &item, err
}

func (i IpPoolController) PostBatch() error {
	var req dto.IpPoolOp
	err := i.Ctx.ReadJSON(&req)
	if err != nil {
		return err
	}
	return i.IpPoolService.Batch(req)
}

func (i IpPoolController) GetBy(name string) (dto.IpPool, error) {
	return i.IpPoolService.Get(name)
}
