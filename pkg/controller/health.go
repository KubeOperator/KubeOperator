package controller

import (
	"fmt"
	"github.com/kataras/iris/v12/context"
)

type Info struct {
	Status string `json:"status"`
	Msg string 	`json:"message"`
}

func HealthController(ctx context.Context) {
	info := Info{
		Status: "1",
		Msg: "Success",
	}
	if _, err := ctx.JSON(info); err != nil {
		fmt.Printf("HealthController err: %v", err)
	}
}

