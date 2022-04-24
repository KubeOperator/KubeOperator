package middleware

import (
	"errors"
	"net/http"

	"github.com/kataras/iris/v12/context"
)

func LanguageMiddleware(ctx context.Context) {
	q := ctx.Request().URL.Query()
	language := q.Get("l")
	if language != "zh-CN" && language != "en-US" {
		errorHandler(ctx, http.StatusBadRequest, errors.New("incorrect request parameter 'l' "))
		return
	}
	ctx.Next()
}
