package controller

import (
	"github.com/KubeOperator/KubeOperator/pkg/controller/page"
	"github.com/KubeOperator/KubeOperator/pkg/service"
	"github.com/kataras/iris/v12/context"
)

type CloudProviderController struct {
	Ctx                  context.Context
	CloudProviderService service.CloudProviderService
}

func NewCloudProviderController() *CloudProviderController {
	return &CloudProviderController{
		CloudProviderService: service.NewCloudProviderService(),
	}
}

func (c CloudProviderController) Get() (page.Page, error) {

	var page page.Page
	items, err := c.CloudProviderService.List()
	if err != nil {
		return page, err
	}
	page.Items = items
	page.Total = len(items)
	return page, nil
}
