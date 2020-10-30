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
	authParty := v1.Party("/auth")
	mvc.New(authParty.Party("/session")).Handle(ErrorHandler).Handle(controller.NewSessionController())
	authScope := v1.Party("/")
	authScope.Use(middleware.JWTMiddleware().Serve)
	authScope.Use(middleware.PagerMiddleware)
	authScope.Use(middleware.UserMiddleware)
	mvc.New(authScope.Party("/clusters")).HandleError(ErrorHandler).Handle(controller.NewClusterController())
	mvc.New(authScope.Party("/credentials")).HandleError(ErrorHandler).Handle(controller.NewCredentialController())
	mvc.New(authScope.Party("/hosts")).HandleError(ErrorHandler).Handle(controller.NewHostController())
	mvc.New(authScope.Party("/users")).HandleError(ErrorHandler).Handle(controller.NewUserController())
	mvc.New(authScope.Party("/regions")).HandleError(ErrorHandler).Handle(controller.NewRegionController())
	mvc.New(authScope.Party("/zones")).HandleError(ErrorHandler).Handle(controller.NewZoneController())
	mvc.New(authScope.Party("/plans")).HandleError(ErrorHandler).Handle(controller.NewPlanController())
	mvc.New(authScope.Party("/systemSettings")).HandleError(ErrorHandler).Handle(controller.NewSystemSettingController())
	mvc.New(authScope.Party("/projects")).HandleError(ErrorHandler).Handle(controller.NewProjectController())
	mvc.New(authScope.Party("/project/resources")).HandleError(ErrorHandler).Handle(controller.NewProjectResourceController())
	mvc.New(authScope.Party("/project/members")).HandleError(ErrorHandler).Handle(controller.NewProjectMemberController())
	mvc.New(authScope.Party("/backupAccounts")).HandleError(ErrorHandler).Handle(controller.NewBackupAccountController())
	mvc.New(authScope.Party("/cluster/backup")).HandleError(ErrorHandler).Handle(controller.NewClusterBackupStrategyController())
	mvc.New(authScope.Party("/license")).Handle(ErrorHandler).Handle(controller.NewLicenseController())
	mvc.New(authScope.Party("/cluster/backup/files")).HandleError(ErrorHandler).Handle(controller.NewClusterBackupFileController())
	mvc.New(authScope.Party("/manifests")).HandleError(ErrorHandler).Handle(controller.NewManifestController())
	mvc.New(authScope.Party("/vm/configs")).HandleError(ErrorHandler).Handle(controller.NewVmConfigController())
	mvc.New(authScope.Party("/events")).HandleError(ErrorHandler).Handle(controller.NewClusterEventController())
	white := v1.Party("/")
	white.Get("/clusters/kubeconfig/{name}", downloadKubeconfig)
	white.Get("/captcha", generateCaptcha)
	mvc.New(white.Party("/theme")).HandleError(ErrorHandler).Handle(controller.NewThemeController())

}

func ErrorHandler(ctx context.Context, err error) {
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
