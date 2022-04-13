package middleware

import (
	"net/http"

	"github.com/KubeOperator/KubeOperator/pkg/service"
	"github.com/KubeOperator/KubeOperator/pkg/util/encrypt"
	"github.com/kataras/iris/v12/context"
)

func HelperMiddleware(ctx context.Context) {
	svc := service.NewLicenseService()
	license, err := svc.GetHw()
	if err != nil || license == nil || encrypt.Md5Str(license.Status) != "9f7d0ee82b6a6ca7ddeae841f3253059" || encrypt.Md5Str(license.Product) != "9d383d901aab38079ba1b7d83403a610" {
		ctx.StatusCode(http.StatusForbidden)
		_, _ = ctx.JSON(map[string]interface{}{
			"msg": "The license is not activated or the authorization has expired!",
		})
		ctx.StopExecution()
		return
	}
	ctx.Next()
}
