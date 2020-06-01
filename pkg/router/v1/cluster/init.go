package cluster

import (
	"github.com/KubeOperator/KubeOperator/pkg/router/v1/cluster/serializer"
	clusterService "github.com/KubeOperator/KubeOperator/pkg/service/cluster"
	"github.com/gin-gonic/gin"
	"net/http"
)

// InitCluster
// @Tags Cluster
// @Summary Cluster
// @Description Init Cluster
// @Accept  json
// @Produce json
// @Param cluster_name path string true "cluster name"
// @Success 200 {object} serializer.InitClusterResponse
// @Router /clusters/initial/{cluster_name}/ [post]
func Init(ctx *gin.Context) {
	clusterName := ctx.Param("name")
	if clusterName == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": invalidClusterNameError.Error(),
		})
		return
	}
	cluster, err := clusterService.Get(clusterName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}
	if err := clusterService.RetryInitCluster(cluster); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, serializer.InitClusterResponse{Message: "cluster initial running"})
}
