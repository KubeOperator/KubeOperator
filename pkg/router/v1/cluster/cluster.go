package cluster

import (
	"errors"
	"github.com/gin-gonic/gin"
	"ko3-gin/pkg/constant"
	clusterModel "ko3-gin/pkg/model/cluster"
	commonModel "ko3-gin/pkg/model/common"
	"ko3-gin/pkg/router/v1/cluster/serializer"
	clusterService "ko3-gin/pkg/service/cluster"
	"net/http"
)

// PageCluster
// @Summary Cluster
// @Description List all clusters with page
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
			_ = ctx.Error(err)
			return
		}
	} else {
		ms, err := clusterService.List()
		models = ms
		total = len(ms)
		if err != nil {
			_ = ctx.Error(err)
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

var invalidClusterName = errors.New("invalid cluster name")

func Get(ctx *gin.Context) {
	name := ctx.Query("name")
	if name == "" {
		_ = ctx.Error(invalidClusterName)
	}
	model, err := clusterService.Get(name)
	if err != nil {
		_ = ctx.Error(err)
	}
	ctx.JSON(http.StatusOK, serializer.GetResponse{Item: serializer.FromModel(model)})

}

func Create(ctx *gin.Context) {
	var req serializer.CreateRequest
	err := ctx.ShouldBind(&req)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	model := clusterModel.Cluster{
		BaseModel: commonModel.BaseModel{
			Name: req.Name,
		},
	}
	err = clusterService.Save(&model)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusCreated, serializer.CreateResponse{Item: serializer.FromModel(model)})
}

func Update(ctx *gin.Context) {
	var req serializer.UpdateRequest
	err := ctx.ShouldBind(&req)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	model := clusterModel.Cluster{
		BaseModel: commonModel.BaseModel{
			Name: req.Name,
		},
	}
	err = clusterService.Save(&model)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, serializer.UpdateResponse{Item: serializer.FromModel(model)})

}

func Delete(ctx *gin.Context) {
	name := ctx.Param("name")
	err := clusterService.Delete(name)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, serializer.DeleteResponse{})
}

func Batch(ctx *gin.Context) {
	var req serializer.BatchRequest
	err := ctx.ShouldBind(&req)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	models := make([]clusterModel.Cluster, 0)
	for _, item := range req.Items {
		models = append(models, serializer.ToModel(item))
	}
	models, err = clusterService.Batch(req.Operation, models)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	var resp serializer.BatchResponse

	for _, model := range models {
		resp.Items = append(resp.Items, serializer.FromModel(model))
	}
	ctx.JSON(http.StatusOK, resp)
}
