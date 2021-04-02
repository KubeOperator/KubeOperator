package controller

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/controller/page"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/service"
	"github.com/kataras/iris/v12/context"
)

type ClusterMemberController struct {
	Ctx                  context.Context
	ClusterMemberService service.ClusterMemberService
}

func NewClusterMemberController() *ClusterMemberController {
	return &ClusterMemberController{
		ClusterMemberService: service.NewClusterMemberService(),
	}
}

func (c ClusterMemberController) Get() (*page.Page, error) {
	clusterName := c.Ctx.Params().GetString("cluster")
	p, _ := c.Ctx.Values().GetBool("page")
	if p {
		num, _ := c.Ctx.Values().GetInt(constant.PageNumQueryKey)
		size, _ := c.Ctx.Values().GetInt(constant.PageSizeQueryKey)
		return c.ClusterMemberService.Page(clusterName, num, size)
	} else {
		var page page.Page
		return &page, nil
	}
}

func (c ClusterMemberController) Post() ([]dto.ClusterMember, error) {
	clusterName := c.Ctx.Params().GetString("cluster")
	var req dto.ClusterMemberCreate
	err := c.Ctx.ReadJSON(&req)
	if err != nil {
		return nil, err
	}
	result, err := c.ClusterMemberService.Create(clusterName, req)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c ClusterMemberController) GetUsers() (dto.UsersResponse, error) {
	name := c.Ctx.URLParam("name")
	return c.ClusterMemberService.GetUsers(name)
}
