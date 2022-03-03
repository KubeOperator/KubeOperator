package controller

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/controller/kolog"
	"github.com/KubeOperator/KubeOperator/pkg/controller/koregexp"
	"github.com/KubeOperator/KubeOperator/pkg/controller/page"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/service"
	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12/context"
)

type RegionController struct {
	Ctx           context.Context
	RegionService service.RegionService
}

func NewRegionController() *RegionController {
	return &RegionController{
		RegionService: service.NewRegionService(),
	}
}

func (r RegionController) Get() (page.Page, error) {

	p, _ := r.Ctx.Values().GetBool("page")
	if p {
		num, _ := r.Ctx.Values().GetInt(constant.PageNumQueryKey)
		size, _ := r.Ctx.Values().GetInt(constant.PageSizeQueryKey)
		return r.RegionService.Page(num, size)
	} else {
		var page page.Page
		items, err := r.RegionService.List()
		if err != nil {
			return page, err
		}
		page.Items = items
		page.Total = len(items)
		return page, nil
	}
}

func (r RegionController) GetBy(name string) (dto.Region, error) {
	return r.RegionService.Get(name)
}

func (r RegionController) Post() (*dto.Region, error) {
	var req dto.RegionCreate
	err := r.Ctx.ReadJSON(&req)
	if err != nil {
		return nil, err
	}
	validate := validator.New()
	if err := validate.RegisterValidation("koname", koregexp.CheckNamePattern); err != nil {
		return nil, err
	}
	if err := validate.Struct(req); err != nil {
		return nil, err
	}

	go kolog.Save(r.Ctx, constant.CREATE_REGION, req.Name)

	return r.RegionService.Create(req)
}

func (r RegionController) Delete(name string) error {
	go kolog.Save(r.Ctx, constant.DELETE_REGION, name)

	return r.RegionService.Delete(name)
}

func (r RegionController) PostBatch() error {
	var req dto.RegionOp
	err := r.Ctx.ReadJSON(&req)
	if err != nil {
		return err
	}
	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		return err
	}
	err = r.RegionService.Batch(req)
	if err != nil {
		return err
	}

	delRegions := ""
	for _, item := range req.Items {
		delRegions += (item.Name + ",")
	}
	go kolog.Save(r.Ctx, constant.DELETE_REGION, delRegions)

	return err
}

func (r RegionController) PostDatacenter() (*dto.CloudRegionResponse, error) {
	var req dto.RegionDatacenterRequest
	err := r.Ctx.ReadJSON(&req)
	if err != nil {
		return nil, err
	}

	data, err := r.RegionService.ListDatacenter(req)
	if err != nil {
		return nil, err
	}
	return &dto.CloudRegionResponse{Result: data}, err
}
