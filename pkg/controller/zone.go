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

type ZoneController struct {
	Ctx         context.Context
	ZoneService service.ZoneService
}

func NewZoneController() *ZoneController {
	return &ZoneController{
		ZoneService: service.NewZoneService(),
	}
}

// List Zone
// @Tags zones
// @Summary Show all zones
// @Description Show zones
// @Accept  json
// @Produce  json
// @Success 200 {object} page.Page
// @Security ApiKeyAuth
// @Router /zones/ [get]
func (z ZoneController) Get() (*page.Page, error) {

	p, _ := z.Ctx.Values().GetBool("page")
	if p {
		num, _ := z.Ctx.Values().GetInt(constant.PageNumQueryKey)
		size, _ := z.Ctx.Values().GetInt(constant.PageSizeQueryKey)
		return z.ZoneService.Page(num, size, condition.TODO())
	} else {
		var page page.Page
		items, err := z.ZoneService.List(condition.TODO())
		if err != nil {
			return nil, err
		}
		page.Items = items
		page.Total = len(items)
		return &page, nil
	}
}

func (z ZoneController) PostSearch() (*page.Page, error) {
	p, _ := z.Ctx.Values().GetBool("page")
	var conditions condition.Conditions
	if z.Ctx.GetContentLength() > 0 {
		if err := z.Ctx.ReadJSON(&conditions); err != nil {
			return nil, err
		}
	}
	if p {
		num, _ := z.Ctx.Values().GetInt(constant.PageNumQueryKey)
		size, _ := z.Ctx.Values().GetInt(constant.PageSizeQueryKey)
		return z.ZoneService.Page(num, size, conditions)

	} else {
		var p page.Page
		items, err := z.ZoneService.List(conditions)
		if err != nil {
			return nil, err
		}
		p.Items = items
		p.Total = len(items)
		return &p, nil
	}
}

// Get Zone
// @Tags zones
// @Summary Show a zone
// @Description show a zone by name
// @Accept  json
// @Produce  json
// @Success 200 {object} dto.Zone
// @Security ApiKeyAuth
// @Router /zones/{name}/ [get]
func (z ZoneController) GetBy(name string) (*dto.Zone, error) {
	return z.ZoneService.Get(name)
}

func (z ZoneController) GetListBy(regionName string) ([]dto.Zone, error) {
	return z.ZoneService.ListByRegionName(regionName)
}

// Create Zone
// @Tags zones
// @Summary Create a zone
// @Description create a zone
// @Accept  json
// @Produce  json
// @Param request body dto.ZoneCreate true "request"
// @Success 200 {object} dto.Zone
// @Security ApiKeyAuth
// @Router /zones/ [post]
func (z ZoneController) Post() (*dto.Zone, error) {
	var req dto.ZoneCreate
	err := z.Ctx.ReadJSON(&req)
	if err != nil {
		return nil, err
	}
	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		return nil, err
	}

	operator := z.Ctx.Values().GetString("operator")
	go kolog.Save(operator, constant.CREATE_ZONE, req.Name)

	return z.ZoneService.Create(req)
}

// Delete Zone
// @Tags zones
// @Summary Delete a zone
// @Description delete a zone by name
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Router /zones/{name}/ [delete]
func (z ZoneController) DeleteBy(name string) error {
	operator := z.Ctx.Values().GetString("operator")
	go kolog.Save(operator, constant.DELETE_ZONE, name)

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

	operator := z.Ctx.Values().GetString("operator")
	go kolog.Save(operator, constant.UPDATE_ZONE, name)

	return z.ZoneService.Update(name, req)
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

	operator := z.Ctx.Values().GetString("operator")
	delZone := ""
	for _, item := range req.Items {
		delZone += (item.Name + ",")
	}
	go kolog.Save(operator, constant.DELETE_ZONE, delZone)

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
