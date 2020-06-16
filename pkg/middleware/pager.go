package middleware

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/kataras/iris/v12/context"
	"strconv"
)

func PagerMiddleware(ctx context.Context) {
	num := ctx.URLParam(constant.PageNumQueryKey)
	limit := ctx.URLParam(constant.PageSizeQueryKey)
	limitInt, err := strconv.Atoi(limit)
	if err != nil || limitInt < 0 {
		ctx.Values().Set("page", false)
		ctx.Next()
	}
	numInt, err := strconv.Atoi(num)
	if err != nil || numInt < 0 {
		ctx.Values().Set("page", false)
		ctx.Next()
	}
	ctx.Values().Set("page", true)
	ctx.Values().Set(constant.PageNumQueryKey, numInt)
	ctx.Values().Set(constant.PageSizeQueryKey, limitInt)
	ctx.Next()
}
