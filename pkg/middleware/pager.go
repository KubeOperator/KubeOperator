package middleware

import (
	"github.com/gin-gonic/gin"
	"ko3-gin/pkg/constant"
	"strconv"
)

func PagerMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		num := ctx.Query(constant.PageNumQueryKey)
		limit := ctx.Query(constant.PageNumQueryKey)
		limitInt, err := strconv.Atoi(limit)
		if err != nil || limitInt < 0 {
			ctx.Set("page", false)
			ctx.Next()
		}
		numInt, err := strconv.Atoi(num)
		if err != nil || numInt < 0 {
			ctx.Set("page", false)
			ctx.Next()
		}
		ctx.Set("page", true)
		ctx.Set(constant.PageNumQueryKey, numInt)
		ctx.Set(constant.PageNumQueryKey, limitInt)
		ctx.Next()
	}
}
