package controller

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/controller/page"
	"github.com/kataras/iris/v12/context"
)

type CloudProviderController struct {
	Ctx context.Context
}

func NewCloudProviderController() *CloudProviderController {
	return &CloudProviderController{}
}

func (c CloudProviderController) Get() (page.Page, error) {

	var page page.Page

	var items []string
	items = append(items, constant.OpenStack)
	items = append(items, constant.VSphere)

	page.Items = items
	page.Total = len(items)
	return page, nil
}
