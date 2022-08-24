package controller

import (
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/service"
	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12/context"
)

type KubernetesController struct {
	Ctx               context.Context
	KubernetesService service.KubernetesService
}

func NewKubernetesController() *KubernetesController {
	return &KubernetesController{
		KubernetesService: service.NewKubernetesService(),
	}
}

func (k KubernetesController) PostSearch() (interface{}, error) {
	var req dto.SourceSearch
	err := k.Ctx.ReadJSON(&req)
	if err != nil {
		return dto.SourceList{}, err
	}
	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		return dto.SourceList{}, err
	}

	return k.KubernetesService.Get(req)
}

func (k KubernetesController) PostSearchMetricBy(cluster string) (interface{}, error) {
	return k.KubernetesService.GetMetric(cluster)
}

func (k KubernetesController) PostCreateSc() error {
	var req dto.SourceScCreate
	if err := k.Ctx.ReadJSON(&req); err != nil {
		return err
	}

	return k.KubernetesService.CreateSc(req)
}

func (k KubernetesController) PostCreateSecret() error {
	var req dto.SourceSecretCreate
	if err := k.Ctx.ReadJSON(&req); err != nil {
		return err
	}

	return k.KubernetesService.CreateSecret(req)
}

func (k KubernetesController) PostCordon() error {
	var req dto.Cordon
	if err := k.Ctx.ReadJSON(&req); err != nil {
		return err
	}

	return k.KubernetesService.CordonNode(req)
}

func (k KubernetesController) PostEvict() error {
	var req dto.Evict
	if err := k.Ctx.ReadJSON(&req); err != nil {
		return err
	}

	return k.KubernetesService.EvictPod(req)
}

func (k KubernetesController) PostDelete() error {
	var req dto.SourceDelete
	if err := k.Ctx.ReadJSON(&req); err != nil {
		return err
	}

	return k.KubernetesService.Delete(req)
}
