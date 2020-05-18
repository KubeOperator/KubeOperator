package cluster

import (
	"github.com/gin-gonic/gin"
	"ko3-gin/pkg/constant"
	clusterModel "ko3-gin/pkg/model/cluster"
	commonModel "ko3-gin/pkg/model/common"
	"ko3-gin/pkg/router/v1/common"
	clusterService "ko3-gin/pkg/service/cluster"
	"net/http"
	"strconv"
)

func List(ctx *gin.Context) {
	num := ctx.Query(constant.PageNum)
	size := ctx.Query(constant.PageSize)
	if num != "" && size != "" {
		pageNum, err := strconv.Atoi(num)
		if err != nil {
			_ = ctx.Error(common.InvalidPageParam)
			return
		}
		pageSize, err := strconv.Atoi(size)
		if err != nil {
			_ = ctx.Error(common.InvalidPageParam)
			return
		}
		models, total, err := clusterService.Page(pageNum, pageSize)
		if err != nil {
			_ = ctx.Error(err)
		}
		items := make([]Cluster, 0)
		for _, model := range models {
			items = append(items, FromModel(model))
		}
		ctx.JSON(http.StatusOK, common.PageResponse{
			Items: items,
			Total: total,
		})
	} else {
		models, err := clusterService.List()
		items := make([]Cluster, 0)
		for _, model := range models {
			items = append(items, FromModel(model))
		}
		if err != nil {
			_ = ctx.Error(err)
		}
		ctx.JSON(http.StatusOK, ListResponse{items: items})
	}
}

func Create(ctx *gin.Context) {
	var req CreateRequest
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
	ctx.JSON(http.StatusCreated, CreateResponse{Item: FromModel(model)})
}

func Update(ctx *gin.Context) {
	var req UpdateRequest
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
	ctx.JSON(http.StatusOK, UpdateResponse{Item: FromModel(model)})

}

func Delete(ctx *gin.Context) {
	var req DeleteRequest
	err := ctx.ShouldBind(&req)
	if err != nil {
		_ = ctx.Error(err)
	}
	err = clusterService.Delete(req.Name)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, DeleteResponse{})
}
