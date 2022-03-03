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

type ZoneController struct {
	Ctx         context.Context
	ZoneService service.ZoneService
}

func NewZoneController() *ZoneController {
	return &ZoneController{
		ZoneService: service.NewZoneService(),
	}
}

func (z ZoneController) Get() (page.Page, error) {

	p, _ := z.Ctx.Values().GetBool("page")
	if p {
		num, _ := z.Ctx.Values().GetInt(constant.PageNumQueryKey)
		size, _ := z.Ctx.Values().GetInt(constant.PageSizeQueryKey)
		return z.ZoneService.Page(num, size)
	} else {
		var page page.Page
		items, err := z.ZoneService.List()
		if err != nil {
			return page, err
		}
		page.Items = items
		page.Total = len(items)
		return page, nil
	}
}

func (z ZoneController) GetBy(name string) (dto.Zone, error) {
	return z.ZoneService.Get(name)
}

func (z ZoneController) GetListBy(regionName string) ([]dto.Zone, error) {
	return z.ZoneService.ListByRegionName(regionName)
}

func (z ZoneController) Post() (*dto.Zone, error) {
	var req dto.ZoneCreate
	err := z.Ctx.ReadJSON(&req)
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

	go kolog.Save(z.Ctx, constant.CREATE_ZONE, req.Name)

	return z.ZoneService.Create(req)
}

func (z ZoneController) Delete(name string) error {
	go kolog.Save(z.Ctx, constant.DELETE_ZONE, name)

	return z.ZoneService.Delete(name)
}

func (z ZoneController) PatchBy(name string) (*dto.Zone, error) {
	var req dto.ZoneUpdate
	err := z.Ctx.ReadJSON(&req)
	if err != nil {
		return nil, err
	}
	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		return nil, err
	}

	go kolog.Save(z.Ctx, constant.UPDATE_ZONE, name)

	return z.ZoneService.Update(req)
}

func (z ZoneController) PostBatch() error {
	var req dto.ZoneOp
	err := z.Ctx.ReadJSON(&req)
	if err != nil {
		return err
	}
	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		return err
	}
	err = z.ZoneService.Batch(req)
	if err != nil {
		return err
	}

	delZone := ""
	for _, item := range req.Items {
		delZone += (item.Name + ",")
	}
	go kolog.Save(z.Ctx, constant.DELETE_ZONE, delZone)

	return err
}

func (z ZoneController) PostClusters() (dto.CloudZoneResponse, error) {
	var req dto.CloudZoneRequest
	err := z.Ctx.ReadJSON(&req)
	if err != nil {
		return dto.CloudZoneResponse{}, err
	}

	data, err := z.ZoneService.ListClusters(req)
	if err != nil {
		return dto.CloudZoneResponse{}, err
	}

	return dto.CloudZoneResponse{Result: data}, err
}

func (z ZoneController) PostTemplates() (dto.CloudZoneResponse, error) {
	var req dto.CloudZoneRequest
	err := z.Ctx.ReadJSON(&req)
	if err != nil {
		return dto.CloudZoneResponse{}, err
	}

	data, err := z.ZoneService.ListTemplates(req)
	if err != nil {
		return dto.CloudZoneResponse{}, err
	}

	return dto.CloudZoneResponse{Result: data}, err
}

func (z ZoneController) PostDatastores() ([]dto.CloudDatastore, error) {
	var req dto.CloudZoneRequest
	err := z.Ctx.ReadJSON(&req)
	if err != nil {
		return nil, err
	}
	return z.ZoneService.ListDatastores(req)
}
