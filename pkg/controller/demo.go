package controller

import (
	"fmt"
	"github.com/kataras/iris"
)

type demoController struct {
	Ctx iris.Context
}

func NewDemoController() *demoController {
	return &demoController{}
}

func (d demoController) GetBy(name string) string {
	fmt.Println(d.Ctx.Path())
	return fmt.Sprintf("hello: %s", name)
}

func (d demoController) Get() []string {
	return []string{"a", "b", "c"}
}
