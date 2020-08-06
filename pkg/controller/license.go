package controller

import (
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/service"
	"github.com/kataras/iris/v12/context"
	"io/ioutil"
)

type LicenseController struct {
	Ctx            context.Context
	LicenseService service.LicenseService
}

func NewLicenseController() *LicenseController {
	return &LicenseController{
		LicenseService: service.NewLicenseService(),
	}
}

func (l *LicenseController) Get() (*dto.License, error) {
	return l.LicenseService.Get()
}
func (l *LicenseController) Post() (*dto.License, error) {
	f, _, err := l.Ctx.FormFile("file")
	if err != nil {
		return nil, err
	}
	bs, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return l.LicenseService.Save(string(bs))
}
