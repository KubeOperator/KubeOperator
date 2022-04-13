package controller

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/controller/kolog"
	"github.com/KubeOperator/KubeOperator/pkg/controller/koregexp"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/service"
	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12/context"
)

type ClusterHelperController struct {
	Ctx            context.Context
	ClusterService service.ClusterService
	LicenseService service.LicenseService
}

func NewClusterHelperController() *ClusterHelperController {
	return &ClusterHelperController{
		ClusterService: service.NewClusterService(),
		LicenseService: service.NewLicenseService(),
	}
}

func (c ClusterHelperController) Post() (*dto.Cluster, error) {
	var req dto.ClusterCreate
	if err := c.Ctx.ReadJSON(&req); err != nil {
		return nil, err
	}

	validate := validator.New()
	if err := validate.RegisterValidation("clustername", koregexp.CheckClusterNamePattern); err != nil {
		return nil, err
	}

	if err := validate.Struct(req); err != nil {
		return nil, err
	}

	item, err := c.ClusterService.Create(req)
	if err != nil {
		return nil, err
	}
	go kolog.Save(c.Ctx, constant.CREATE_CLUSTER, req.Name)
	return item, nil
}
