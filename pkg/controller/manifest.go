package controller

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/controller/log"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/service"
	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12/context"
)

type ManifestController struct {
	Ctx             context.Context
	ManifestService service.ClusterManifestService
}

func NewManifestController() *ManifestController {
	return &ManifestController{
		ManifestService: service.NewClusterManifestService(),
	}
}

func (m *ManifestController) Get() ([]dto.ClusterManifest, error) {
	return m.ManifestService.List()
}

func (m *ManifestController) GetActive() ([]dto.ClusterManifest, error) {
	return m.ManifestService.ListActive()
}

func (m ManifestController) PatchBy(name string) (model.ClusterManifest, error) {
	var req dto.ClusterManifestUpdate
	err := m.Ctx.ReadJSON(&req)

	if err != nil {
		return model.ClusterManifest{}, err
	}
	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		return model.ClusterManifest{}, err
	}

	operator := m.Ctx.Values().GetString("operator")
	if req.IsActive {
		go log.Save(operator, constant.ENABLE_VERSION, req.Name)
	} else {
		go log.Save(operator, constant.DISABLE_VERSION, req.Name)
	}

	return m.ManifestService.Update(req)
}
