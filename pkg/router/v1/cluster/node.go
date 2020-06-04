package cluster

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	clusterModel "github.com/KubeOperator/KubeOperator/pkg/model/cluster"
	"github.com/KubeOperator/KubeOperator/pkg/router/v1/cluster/serializer"
	clusterService "github.com/KubeOperator/KubeOperator/pkg/service/cluster"
	"github.com/gin-gonic/gin"
	"net/http"
)

// ListClusterNodes
// @Tags ClusterNode
// @Summary ClusterNode
// @Description List ClusterNodes
// @Accept  json
// @Produce json
// @Param pageNum query string false "page num"
// @Param pageSize query string false "page size"
// @Success 200 {object} serializer.ListNodeResponse
// @Router /clusters/{cluster_name}/nodes/ [get]
func ListNodes(ctx *gin.Context) {
	page := ctx.GetBool("page")
	var models []clusterModel.Node
	total := 0
	clusterName := ctx.Param("name")
	if page {
		pageNum := ctx.GetInt(constant.PageNumQueryKey)
		pageSize := ctx.GetInt(constant.PageSizeQueryKey)
		m, t, err := clusterService.PageClusterNodes(clusterName, pageNum, pageSize)
		models = m
		total = t
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"msg": err.Error(),
			})
			return
		}
	} else {
		ms, err := clusterService.ListClusterNodes(clusterName)
		models = ms
		total = len(ms)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"msg": err.Error(),
			})
			return
		}
	}
	var resp = serializer.ListNodeResponse{
		Items: []serializer.Node{},
		Total: total,
	}
	for _, model := range models {
		resp.Items = append(resp.Items, serializer.FromNodeModel(model))
	}
	ctx.JSON(http.StatusOK, resp)
}
