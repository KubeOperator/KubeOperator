package controller

import (
	"github.com/KubeOperator/KubeOperator/pkg/service"
	"github.com/KubeOperator/KubeOperator/pkg/service/dto"
	"github.com/kataras/iris/v12"
)

type HostController struct {
	ctx         iris.Context
	hostService service.HostService
}

func (h HostController) Get() ([]dto.Host, error) {
	return h.hostService.List()
}

func (h HostController) GetBy(name string) (dto.Host, error) {
	return h.hostService.Get(name)
}

func (h HostController) Post() error {
	var req dto.HostCreate
	err := h.ctx.ReadJSON(&req)
	if err != nil {
		return err
	}
	return nil
}

func (h HostController) Delete(name string) error {
	return h.hostService.Delete(name)
}
