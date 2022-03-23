package v1

import (
	"net/http"

	"github.com/KubeOperator/KubeOperator/pkg/util/captcha"
	"github.com/kataras/iris/v12/context"
)

func generateCaptcha(ctx context.Context) {
	c, err := captcha.CreateCaptcha()
	if err != nil {
		_, _ = ctx.JSON(err)
		ctx.StatusCode(http.StatusInternalServerError)
	}
	_, _ = ctx.JSON(&c)

}
