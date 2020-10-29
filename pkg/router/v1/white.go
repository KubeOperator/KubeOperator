package v1

import (
	"fmt"
	"github.com/KubeOperator/KubeOperator/pkg/service"
	"github.com/KubeOperator/KubeOperator/pkg/util/captcha"
	"github.com/kataras/iris/v12/context"
	"net/http"
)

func downloadKubeconfig(ctx context.Context) {
	clusterName := ctx.Params().GetString("name")
	ctx.Header("Content-Disposition", "attachment")
	ctx.Header("filename", fmt.Sprintf("%s-config", clusterName))
	ctx.Header("Content-Type", "application/download")
	clusterService := service.NewClusterService()
	str, err := clusterService.GetKubeconfig(clusterName)
	if err != nil {
		_, _ = ctx.JSON(err)
		ctx.StatusCode(http.StatusInternalServerError)
	}
	_, _ = ctx.WriteString(str)
}

func generateCaptcha(ctx context.Context) {
	c, err := captcha.CreateCaptcha()
	if err != nil {
		_, _ = ctx.JSON(err)
		ctx.StatusCode(http.StatusInternalServerError)
	}
	_, _ = ctx.JSON(&c)

}
