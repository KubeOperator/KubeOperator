package host

import (
	"github.com/gin-gonic/gin"
	"ko3-gin/internal/constant"
	"net/http"
	"strconv"
)

func Create(ctx *gin.Context) {
	var request Host
	err := ctx.Bind(&request)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	if err := Service.Create(Host{
		Name: request.Name,
		Ip:   request.Ip,
		Port: request.Port,
	}); err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, nil)
}

func List(ctx *gin.Context) {
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
		items, total := Service.Page(pageNum, pageSize)
		ctx.JSON(http.StatusOK, PageHostResponse{
			Items: items,
			Total: total,
		})
	}

}
