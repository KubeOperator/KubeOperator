package v1

import (
	"encoding/json"
	"github.com/KubeOperator/KubeOperator/pkg/controller"
	"github.com/jinzhu/gorm"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/mvc"
	"net/http"
)

func V1(parent iris.Party) {
	v1 := parent.Party("/v1")
	mvc.New(v1.Party("/clusters")).HandleError(errorHandler).Handle(controller.NewClusterController())
	mvc.New(v1.Party("/credentials")).HandleError(errorHandler).Handle(controller.NewCredentialController())
	mvc.New(v1.Party("/hosts")).HandleError(errorHandler).Handle(controller.NewHostController())
	mvc.New(v1.Party("/users")).HandleError(errorHandler).Handle(controller.NewUserController())
	mvc.New(v1.Party("/regions")).HandleError(errorHandler).Handle(controller.NewRegionController())
	mvc.New(v1.Party("/cloud/providers")).HandleError(errorHandler).Handle(controller.NewCloudProviderController())
	mvc.New(v1.Party("/zones")).HandleError(errorHandler).Handle(controller.NewZoneController())
	mvc.New(v1.Party("/plans")).HandleError(errorHandler).Handle(controller.NewPlanController())
	mvc.New(v1.Party("/systemSettings")).HandleError(errorHandler).Handle(controller.NewSystemSettingController())
	mvc.New(v1.Party("/projects")).HandleError(errorHandler).Handle(controller.NewProjectController())
	mvc.New(v1.Party("/project/resources")).HandleError(errorHandler).Handle(controller.NewProjectResourceController())
	mvc.New(v1.Party("/project/members")).HandleError(errorHandler).Handle(controller.NewProjectMemberController())
	mvc.New(v1.Party("/backupAccounts")).HandleError(errorHandler).Handle(controller.NewBackupAccountController())
	mvc.New(v1.Party("/cluster/backup")).HandleError(errorHandler).Handle(controller.NewClusterBackupStrategyController())
	mvc.New(v1.Party("/license")).Handle(errorHandler).Handle(controller.NewLicenseController())
	mvc.New(v1.Party("/cluster/backup/files")).Handle(errorHandler).Handle(controller.NewClusterBackupFileController())
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
