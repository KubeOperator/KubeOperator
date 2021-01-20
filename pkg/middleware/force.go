package middleware

import (
	"github.com/kataras/iris/v12/context"
	"strings"
)

func ForceMiddleware(ctx context.Context) {
	f := ctx.URLParam("force")
	if strings.ToLower(f) == "true" {
		ctx.Values().Set("force", true)
	} else {
		ctx.Values().Set("force", false)
	}
	ctx.Next()
}
