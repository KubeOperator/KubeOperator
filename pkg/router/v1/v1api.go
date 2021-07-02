package v1

import (
	"encoding/json"
	"net/http"

	"github.com/KubeOperator/KubeOperator/pkg/controller"
	"github.com/KubeOperator/KubeOperator/pkg/errorf"
	"github.com/KubeOperator/KubeOperator/pkg/middleware"
	"github.com/jinzhu/gorm"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/mvc"
	"github.com/pkg/errors"
)

var AuthScope iris.Party
var WhiteScope iris.Party

func V1(parent iris.Party) {
	v1 := parent.Party("/v1")
	authParty := v1.Party("/auth")
	mvc.New(authParty.Party("/session")).HandleError(ErrorHandler).Handle(controller.NewSessionController())
	mvc.New(v1.Party("/user")).HandleError(ErrorHandler).Handle(controller.NewForgotPasswordController())
	AuthScope = v1.Party("/")
	AuthScope.Use(middleware.JWTMiddleware().Serve)
	AuthScope.Use(middleware.UserMiddleware)
	AuthScope.Use(middleware.RBACMiddleware())
	AuthScope.Use(middleware.PagerMiddleware)
	AuthScope.Use(middleware.ForceMiddleware)
	mvc.New(AuthScope.Party("/clusters")).HandleError(ErrorHandler).Handle(controller.NewClusterController())
	mvc.New(AuthScope.Party("/credentials")).HandleError(ErrorHandler).Handle(controller.NewCredentialController())
	mvc.New(AuthScope.Party("/hosts")).HandleError(ErrorHandler).Handle(controller.NewHostController())
	mvc.New(AuthScope.Party("/users")).HandleError(ErrorHandler).Handle(controller.NewUserController())
	mvc.New(AuthScope.Party("/regions")).HandleError(ErrorHandler).Handle(controller.NewRegionController())
	mvc.New(AuthScope.Party("/zones")).HandleError(ErrorHandler).Handle(controller.NewZoneController())
	mvc.New(AuthScope.Party("/plans")).HandleError(ErrorHandler).Handle(controller.NewPlanController())
	mvc.New(AuthScope.Party("/settings")).HandleError(ErrorHandler).Handle(controller.NewSystemSettingController())
	mvc.New(AuthScope.Party("/logs")).HandleError(ErrorHandler).Handle(controller.NewSystemLogController())
	mvc.New(AuthScope.Party("/projects")).HandleError(ErrorHandler).Handle(controller.NewProjectController())
	mvc.New(AuthScope.Party("/clusters/istio")).HandleError(ErrorHandler).Handle(controller.NewClusterIstioController())
	mvc.New(AuthScope.Party("/backupaccounts")).HandleError(ErrorHandler).Handle(controller.NewBackupAccountController())
	mvc.New(AuthScope.Party("/clusters/backup")).HandleError(ErrorHandler).Handle(controller.NewClusterBackupStrategyController())
	mvc.New(AuthScope.Party("/license")).Handle(ErrorHandler).Handle(controller.NewLicenseController())
	mvc.New(AuthScope.Party("/clusters/backup/files")).HandleError(ErrorHandler).Handle(controller.NewClusterBackupFileController())
	mvc.New(AuthScope.Party("/manifests")).HandleError(ErrorHandler).Handle(controller.NewManifestController())
	mvc.New(AuthScope.Party("/vmconfigs")).HandleError(ErrorHandler).Handle(controller.NewVmConfigController())
	mvc.New(AuthScope.Party("/clusters/events")).HandleError(ErrorHandler).Handle(controller.NewClusterEventController())
	mvc.New(AuthScope.Party("/ippools")).HandleError(ErrorHandler).Handle(controller.NewIpPoolController())
	mvc.New(AuthScope.Party("/ippools/{name}/ips")).HandleError(ErrorHandler).Handle(controller.NewIpController())
	mvc.New(AuthScope.Party("/projects/{project}/resources")).HandleError(ErrorHandler).Handle(controller.NewProjectResourceController())
	mvc.New(AuthScope.Party("/projects/{project}/members")).HandleError(ErrorHandler).Handle(controller.NewProjectMemberController())
	mvc.New(AuthScope.Party("/projects/{project}/clusters/{cluster}/members")).HandleError(ErrorHandler).Handle(controller.NewClusterMemberController())
	mvc.New(AuthScope.Party("/projects/{project}/clusters/{cluster}/resources")).HandleError(ErrorHandler).Handle(controller.NewClusterResourceController())
	WhiteScope = v1.Party("/")
	WhiteScope.Get("/clusters/kubeconfig/{name}", downloadKubeconfig)
	WhiteScope.Get("/captcha", generateCaptcha)
	mvc.New(WhiteScope.Party("/theme")).HandleError(ErrorHandler).Handle(controller.NewThemeController())

}

func ErrorHandler(ctx context.Context, err error) {
	if err != nil {
		warp := struct {
			Msg string `json:"msg"`
		}{err.Error()}
		var result string
		switch errType := err.(type) {
		case gorm.Errors:
			errorSet := make(map[string]string)
			for _, er := range errType {
				tr := ctx.Tr(er.Error())
				if tr != "" {
					errorMsg := tr
					errorSet[er.Error()] = errorMsg
				}
			}
			for _, set := range errorSet {
				result = result + set + " "
			}
		case error:
			switch errRoot := errors.Cause(err).(type) {
			case errorf.CErrFs:
				errs := errRoot.Get()
				for _, er := range errs {
					args := er.Args.([]interface{})
					tr := ctx.Tr(er.Msg, args...)
					if tr != "" {
						result = result + tr + "\n "
					}
				}
			default:
				tr := ctx.Tr(errors.Cause(err).Error())
				if tr != "" {
					result = tr
				} else {
					result = err.Error()
				}
			}
		}
		warp.Msg = result
		bf, _ := json.Marshal(&warp)
		ctx.StatusCode(http.StatusBadRequest)
		_, _ = ctx.WriteString(string(bf))
		ctx.StopExecution()
		return
	}
}
