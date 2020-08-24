package v1

import (
	"encoding/json"
	"github.com/KubeOperator/KubeOperator/pkg/controller"
	"github.com/KubeOperator/KubeOperator/pkg/middleware"
	"github.com/jinzhu/gorm"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/mvc"
	"net/http"
)

func V1(parent iris.Party) {
	v1 := parent.Party("/v1")
	auth := v1.Party("/")
	auth.Use(middleware.PagerMiddleware)
	auth.Use(middleware.JWTMiddleware().Serve)
	auth.Use(middleware.UserMiddleware)
	mvc.New(auth.Party("/clusters")).HandleError(errorHandler).Handle(controller.NewClusterController())
	mvc.New(auth.Party("/credentials")).HandleError(errorHandler).Handle(controller.NewCredentialController())
	mvc.New(auth.Party("/hosts")).HandleError(errorHandler).Handle(controller.NewHostController())
	mvc.New(auth.Party("/users")).HandleError(errorHandler).Handle(controller.NewUserController())
	mvc.New(auth.Party("/regions")).HandleError(errorHandler).Handle(controller.NewRegionController())
	mvc.New(auth.Party("/cloud/providers")).HandleError(errorHandler).Handle(controller.NewCloudProviderController())
	mvc.New(auth.Party("/zones")).HandleError(errorHandler).Handle(controller.NewZoneController())
	mvc.New(auth.Party("/plans")).HandleError(errorHandler).Handle(controller.NewPlanController())
	mvc.New(auth.Party("/systemSettings")).HandleError(errorHandler).Handle(controller.NewSystemSettingController())
	mvc.New(auth.Party("/projects")).HandleError(errorHandler).Handle(controller.NewProjectController())
	mvc.New(auth.Party("/project/resources")).HandleError(errorHandler).Handle(controller.NewProjectResourceController())
	mvc.New(auth.Party("/project/members")).HandleError(errorHandler).Handle(controller.NewProjectMemberController())
	mvc.New(auth.Party("/backupAccounts")).HandleError(errorHandler).Handle(controller.NewBackupAccountController())
	mvc.New(auth.Party("/cluster/backup")).HandleError(errorHandler).Handle(controller.NewClusterBackupStrategyController())
	mvc.New(auth.Party("/license")).Handle(errorHandler).Handle(controller.NewLicenseController())
	mvc.New(auth.Party("/cluster/backup/files")).HandleError(errorHandler).Handle(controller.NewClusterBackupFileController())
	white := v1.Party("/")
	white.Get("/clusters/kubeconfig/{name}", downloadKubeconfig)

}

func errorHandler(ctx context.Context, err error) {
	if err != nil {
		warp := struct {
			Msg string `json:"msg"`
		}{err.Error()}
		var result string
		switch err.(type) {
		case gorm.Errors:
			errorSet := make(map[string]string)
			errors, ok := err.(gorm.Errors)
			if ok {
				for _, er := range errors {
					tr := ctx.Tr(er.Error())
					if tr != "" {
						errorMsg := tr
						errorSet[er.Error()] = errorMsg
					}
				}
				for _, set := range errorSet {
					result = result + set + " "
				}
			}
		case error:
			tr := ctx.Tr(err.Error())
			if tr != "" {
				result = tr
			} else {
				result = err.Error()
			}
			break
		}
		warp.Msg = result
		bf, _ := json.Marshal(&warp)
		ctx.StatusCode(http.StatusBadRequest)
		_, _ = ctx.WriteString(string(bf))
		ctx.StopExecution()
		return
	}
}
