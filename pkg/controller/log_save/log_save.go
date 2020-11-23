package log_save

import (
	"fmt"

	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/service"
)

func LogSave(name, operation, operationInfo string) {
	lS := service.NewSystemLogService()
	logInfo := dto.SystemLogCreate{
		Name:          name,
		Operation:     operation,
		OperationInfo: operationInfo,
	}
	err := lS.Create(logInfo)
	if err != nil {
		fmt.Printf("save system logs err, err: %v\n", err)
	}
}
