package kolog

import (
	"fmt"

	"github.com/kataras/iris/v12/context"

	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/logger"
	"github.com/KubeOperator/KubeOperator/pkg/service"
)

var log = logger.Default

func Save(ctx context.Context, operation, operationInfo string) {
	lS := service.NewSystemLogService()
	operator := ctx.Values().GetString("operator")
	ip := ctx.Values().GetString("ipfrom")
	fmt.Println(ip)
	logInfo := dto.SystemLogCreate{
		Name:          operator,
		Operation:     operation,
		OperationInfo: operationInfo,
		IP:            ip,
	}
	if err := lS.Create(logInfo); err != nil {
		log.Errorf("save system logs failed, error: %s", err.Error())
	}
}
