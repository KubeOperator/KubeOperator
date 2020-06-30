package controller

import (
	"github.com/KubeOperator/KubeOperator/pkg/service"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/kataras/iris"
	"log"
)

type demoController struct {
	Ctx         iris.Context
	demoService service.DemoService
}

func NewDemoController() *demoController {
	return &demoController{
		demoService: service.NewDemoService(),
	}
}

func (d demoController) GetBy(name string) dto.Demo {
	dm, err := d.demoService.Get(name)
	if err != nil {
		log.Println(err.Error())
	}
	return dm
}

func (d demoController) Get() []dto.Demo {
	dms, err := d.demoService.List()
	if err != nil {
		log.Println(err.Error())
	}
	return dms
}

func (d demoController) Post() {
	var req dto.CreateDemo
	err := d.Ctx.ReadJSON(&req)
	if err != nil {
		log.Println(err.Error())
	}
	err = d.demoService.Save(req)
	if err != nil {
		log.Println(err.Error())
	}
}
