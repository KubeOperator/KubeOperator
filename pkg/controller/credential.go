package controller

import (
	"github.com/KubeOperator/KubeOperator/pkg/service"
	"github.com/KubeOperator/KubeOperator/pkg/service/dto"
	"github.com/kataras/iris/v12"
)

type credentialController struct {
	ctx               iris.Context
	credentialService service.CredentialService
}

func (c credentialController) Get() ([]dto.Credential, error) {
	return c.credentialService.List()
}

func (c credentialController) GetBy(name string) (dto.Credential, error) {
	return c.credentialService.Get(name)
}

func (c credentialController) Post() error {
	var req dto.CredentialCreate
	err := c.ctx.ReadJSON(&req)
	if err != nil {
		return err
	}
	return nil
}

func (c credentialController) Delete(name string) error {
	return c.credentialService.Delete(name)
}

func (c credentialController) Batch(operation string, items []dto.Credential) error {
	_, err := c.credentialService.Batch(operation, items)
	return err
}
