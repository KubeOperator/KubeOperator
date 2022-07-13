package controller

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/controller/kolog"
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

// List clusterMember
// @Tags clusterMembers
// @Summary Show all clusterMembers
// @Description 获取集群成员列表
// @Accept  json
// @Produce  json
// @Param project path string true "项目名称"
// @Param cluster path string true "集群名称"
// @Success 200 {object} page.Page
// @Security ApiKeyAuth
// @Router /projects/{project}/clusters/{cluster}/members [get]
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

// Create CLusterMember
// @Tags clusterMembers
// @Summary Create a cLusterMember
// @Description 授权成员到集群
// @Accept  json
// @Produce  json
// @Param request body dto.ClusterMemberCreate true "request"
// @Param project path string true "项目名称"
// @Param cluster path string true "集群名称"
// @Success 200 {Array} []dto.CLusterMember
// @Security ApiKeyAuth
// @Router /projects/{project}/clusters/{cluster}/members [post]
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

	operator := c.Ctx.Values().GetString("operator")
	users := ""
	for _, u := range req.Usernames {
		users += (u + ",")
	}
	go kolog.Save(operator, constant.BIND_CLUSTER_MEMBER, users)

	return result, nil
}

func (c ClusterMemberController) GetUsers() (dto.UsersResponse, error) {
	name := c.Ctx.URLParam("name")
	return c.ClusterMemberService.GetUsers(name)
}

func (c ClusterMemberController) GetSearch() (dto.UsersAddResponse, error) {
	clusterName := c.Ctx.Params().GetString("cluster")
	name := c.Ctx.URLParam("name")
	return c.ClusterMemberService.GetUsersByName(clusterName, name)
}

// Delete CLusterMember
// @Tags clusterMembers
// @Summary Delete cLusterMember
// @Description 取消集群人员授权
// @Accept  json
// @Produce  json
// @Param project path string true "项目名称"
// @Param cluster path string true "集群名称"
// @Param name path string true "人员名称"
// @Security ApiKeyAuth
// @Router /projects/{project}/clusters/{cluster}/members/{name} [delete]
func (c ClusterMemberController) DeleteBy(name string) error {
	clusterName := c.Ctx.Params().GetString("cluster")

	operator := c.Ctx.Values().GetString("operator")
	go kolog.Save(operator, constant.UNBIND_CLUSTER_MEMBER, name)

	return c.ClusterMemberService.Delete(name, clusterName)
}
