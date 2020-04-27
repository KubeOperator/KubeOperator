package api

import (
	"github.com/gin-gonic/gin"
	"ko3-gin/internal/constant"
	"ko3-gin/pkg/api/serializer"
	"ko3-gin/pkg/host"
	"ko3-gin/pkg/model"
	"net/http"
	"strconv"
)


// CreateHost
// @Summary Create a Host
// @Description Create a Host
// @Accept  json
// @Produce json
// @Param request body serializer.CreatHostRequest true "create host request"
// @Router /hosts/ [post]
func HostCreate(ctx *gin.Context) {
	var request model.Host
	err := ctx.Bind(&request)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	if err := host.Service.Create(model.Host{
		Name: request.Name,
		Ip:   request.Ip,
		Port: request.Port,
	}); err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, nil)
}

func HostList(ctx *gin.Context) {
	num := ctx.Query(constant.PageNum)
	size := ctx.Query(constant.PageSize)
	if num != "" && size != "" {
		pageNum, err := strconv.Atoi(num)
		if err != nil {
			_ = ctx.Error(err)
			return
		}
		pageSize, err := strconv.Atoi(size)
		if err != nil {
			_ = ctx.Error(err)
			return
		}
		items, total := host.Service.Page(pageNum, pageSize)
		ctx.JSON(http.StatusOK, serializer.PageHostResponse{
			Items: items,
			Total: total,
		})
	}

}
