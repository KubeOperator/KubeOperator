package controller

import (
	"github.com/kataras/iris/v12/context"
)

type Info struct {
	Status string `json:"status"`
	Msg    string `json:"message"`
}

func HealthController(ctx context.Context) {
	info := Info{
		Status: "1",
		Msg:    "Success",
	}
	_, _ = ctx.JSON(info)
}
