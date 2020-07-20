package v1

import (
	"github.com/KubeOperator/KubeOperator/pkg/controller"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

func V1(parent iris.Party) {
	v1 := parent.Party("/v1")
	mvc.New(v1.Party("/clusters")).Handle(controller.NewClusterController())
	mvc.New(v1.Party("/credentials")).Handle(controller.NewCredentialController())
	mvc.New(v1.Party("/hosts")).Handle(controller.NewHostController())
	mvc.New(v1.Party("/users")).Handle(controller.NewUserController())
	mvc.New(v1.Party("/regions")).Handle(controller.NewRegionController())
	mvc.New(v1.Party("/cloud/providers")).Handle(controller.NewCloudProviderController())
	mvc.New(v1.Party("/zones")).Handle(controller.NewZoneController())
	mvc.New(v1.Party("/plans")).Handle(controller.NewPlanController())
	mvc.New(v1.Party("/systemSettings")).Handle(controller.NewSystemSettingController())
	mvc.New(v1.Party("/projects")).Handle(controller.NewProjectController())
}
