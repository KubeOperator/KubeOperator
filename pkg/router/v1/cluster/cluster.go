package cluster

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	clusterModel "github.com/KubeOperator/KubeOperator/pkg/model/cluster"
	commonModel "github.com/KubeOperator/KubeOperator/pkg/model/common"
	"github.com/KubeOperator/KubeOperator/pkg/router/v1/cluster/serializer"
	clusterService "github.com/KubeOperator/KubeOperator/pkg/service/cluster"
	"net/http"
)

var (
	invalidClusterNameError = errors.New("invalid cluster name")
)

// ListCluster
// @Summary Cluster
// @Description List clusters
// @Accept  json
// @Produce json
// @Param pageNum query string false "page num"
// @Param pageSize query string false "page size"
// @Success 200 {object} serializer.ListResponse
// @Router /clusters/ [get]
func List(ctx *gin.Context) {
	page := ctx.GetBool("page")
	var models []clusterModel.Cluster
	total := 0
	if page {
		pageNum := ctx.GetInt(constant.PageNumQueryKey)
		pageSize := ctx.GetInt(constant.PageSizeQueryKey)
		m, t, err := clusterService.Page(pageNum, pageSize)
		models = m
		total = t
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"msg": err.Error(),
			})
			return
		}
	} else {
		ms, err := clusterService.List()
		models = ms
		total = len(ms)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"msg": err.Error(),
			})
			return
		}
	}
	var resp = serializer.ListResponse{
		Items: []serializer.Cluster{},
		Total: total,
	}
	for _, model := range models {
		resp.Items = append(resp.Items, serializer.FromModel(model))
	}
	ctx.JSON(http.StatusOK, resp)
}

// GetCluster
// @Summary Cluster
// @Description Get Cluster
// @Accept  json
// @Produce json
// @Param cluster_name path string true "cluster name"
// @Success 200 {object} serializer.GetResponse
// @Router /clusters/{cluster_name} [get]
func Get(ctx *gin.Context) {
	name := ctx.Param("name")
	if name == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": invalidClusterNameError,
		})
		return
	}
	model, err := clusterService.Get(name)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, serializer.GetResponse{Item: serializer.FromModel(*model)})
}

// CreateCluster
// @Summary Cluster
// @Description Create a Cluster
// @Accept  json
// @Produce json
// @Param request body serializer.CreateRequest true "cluster"
// @Success 201 {object} serializer.CreateResponse
// @Router /clusters/ [post]
func Create(ctx *gin.Context) {
	var req serializer.CreateRequest
	err := ctx.ShouldBind(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
		return
	}
	model := clusterModel.Cluster{
		BaseModel: commonModel.BaseModel{
			Name: req.Name,
		},
	}
	err = clusterService.Save(&model)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusCreated, serializer.CreateResponse{Item: serializer.FromModel(model)})
}

// UpdateCluster
// @Summary Cluster
// @Description Update a Cluster
// @Accept  json
// @Produce json
// @Param request body serializer.UpdateRequest true "cluster"
// @Param cluster_name path string true "cluster name"
// @Success 200 {object} serializer.UpdateResponse
// @Router /clusters/{cluster_name} [patch]
func Update(ctx *gin.Context) {
	var req serializer.UpdateRequest
	err := ctx.ShouldBind(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
		return
	}
	model := serializer.ToModel(req.Item)
	err = clusterService.Save(&model)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, serializer.UpdateResponse{Item: serializer.FromModel(model)})

}

// DeleteCluster
// @Summary Cluster
// @Description Delete a Cluster
// @Accept  json
// @Produce json
// @Param cluster_name path string true "cluster name"
// @Success 200 {object} serializer.DeleteResponse
// @Router /clusters/{cluster_name} [delete]
func Delete(ctx *gin.Context) {
	name := ctx.Param("name")
	if name == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": invalidClusterNameError.Error(),
		})
		return
	}
	err := clusterService.Delete(name)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, serializer.DeleteResponse{})
}

// BatchCluster
// @Summary Cluster
// @Description Batch Clusters
// @Accept  json
// @Produce json
// @Param request body serializer.BatchRequest true "Batch"
// @Success 200 {object} serializer.BatchResponse
// @Router /clusters/batch/ [post]
func Batch(ctx *gin.Context) {
	var req serializer.BatchRequest
	err := ctx.ShouldBind(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
		return
	}
	models := make([]clusterModel.Cluster, 0)
	for _, item := range req.Items {
		models = append(models, serializer.ToModel(item))
	}
	models, err = clusterService.Batch(req.Operation, models)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}
	var resp serializer.BatchResponse
	for _, model := range models {
		resp.Items = append(resp.Items, serializer.FromModel(model))
	}
	ctx.JSON(http.StatusOK, resp)
}
