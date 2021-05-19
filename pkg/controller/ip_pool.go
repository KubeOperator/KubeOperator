package controller

import (
	"errors"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/controller/condition"
	"github.com/KubeOperator/KubeOperator/pkg/controller/kolog"
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

// List IpPool
// @Tags ippools
// @Summary Show all ippools
// @Description 获取IP池列表
// @Accept  json
// @Produce  json
// @Success 200 {object} page.Page
// @Security ApiKeyAuth
// @Router /ippools [get]
func (i IpPoolController) Get() (*page.Page, error) {
	p, _ := i.Ctx.Values().GetBool("page")
	if p {
		num, _ := i.Ctx.Values().GetInt(constant.PageNumQueryKey)
		size, _ := i.Ctx.Values().GetInt(constant.PageSizeQueryKey)
		return i.IpPoolService.Page(num, size, condition.TODO())
	} else {
		var p page.Page
		items, err := i.IpPoolService.List(condition.TODO())
		if err != nil {
			return nil, err
		}
		p.Items = items
		p.Total = len(items)
		return &p, nil
	}
}

// Search IpPool
// @Tags ippools
// @Summary Search IpPool
// @Description 过滤IP池
// @Accept  json
// @Produce  json
// @Param conditions body condition.Conditions true "conditions"
// @Success 200 {object} page.Page
// @Security ApiKeyAuth
// @Router /ippools/search [post]
func (i IpPoolController) PostSearch() (*page.Page, error) {
	p, _ := i.Ctx.Values().GetBool("page")
	var conditions condition.Conditions
	if i.Ctx.GetContentLength() > 0 {
		if err := i.Ctx.ReadJSON(&conditions); err != nil {
			return nil, err
		}
	}
	if p {
		num, _ := i.Ctx.Values().GetInt(constant.PageNumQueryKey)
		size, _ := i.Ctx.Values().GetInt(constant.PageSizeQueryKey)
		return i.IpPoolService.Page(num, size, conditions)
	} else {
		var p page.Page
		items, err := i.IpPoolService.List(conditions)
		if err != nil {
			return nil, err
		}
		p.Items = items
		p.Total = len(items)
		return &p, nil
	}
}

// Create IpPool
// @Tags ippools
// @Summary Create a IpPool
// @Description 创建IP池
// @Accept  json
// @Produce  json
// @Param request body dto.IpPoolCreate true "request"
// @Success 200 {object} dto.IpPool
// @Security ApiKeyAuth
// @Router /ippools [post]
func (i IpPoolController) Post() (*dto.IpPool, error) {
	var req dto.IpPoolCreate
	err := i.Ctx.ReadJSON(&req)
	if err != nil {
		return nil, err
	}
	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		return nil, err
	}
	item, _ := i.IpPoolService.Get(req.Name)
	if item.ID != "" {
		return nil, errors.New("NAME_EXISTS")
	}
	operator := i.Ctx.Values().GetString("operator")
	go kolog.Save(operator, constant.CREATE_IP_POOL, "")

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
	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		return err
	}

	operator := i.Ctx.Values().GetString("operator")
	ips := ""
	for _, ip := range req.Items {
		ips += (ip.Name + ",")
	}
	go kolog.Save(operator, constant.BACTH_DELETE_IP_POOL, ips)

	return i.IpPoolService.Batch(req)
}

// Get IpPool
// @Tags ippools
// @Summary Get IpPool
// @Description 获取单个IP池
// @Accept  json
// @Produce  json
// @Param name path string true "IP池名称"
// @Success 200 {object} dto.IpPool
// @Security ApiKeyAuth
// @Router /ippools/{name} [get]
func (i IpPoolController) GetBy(name string) (dto.IpPool, error) {
	return i.IpPoolService.Get(name)
}

// Delete IpPool
// @Tags ippools
// @Summary Delete a IpPool
// @Description  删除IP池
// @Accept  json
// @Produce  json
// @Param name path string true "IP池名称"
// @Security ApiKeyAuth
// @Router /ippools/{name} [delete]
func (i IpPoolController) DeleteBy(name string) error {
	operator := i.Ctx.Values().GetString("operator")
	go kolog.Save(operator, constant.DELETE_IP_POOL, name)

	return i.IpPoolService.Delete(name)
}
