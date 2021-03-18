package dto

import (
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/util/typeparse"
)

type SystemLog struct {
	model.SystemLog
}

type SystemLogCreate struct {
	Name          string `json:"name"`
	Operation     string `json:"operation"`
	OperationInfo string `json:"operationInfo"`
}

type SystemLogQuery struct {
	Name          typeparse.QueryCondition `json:"name"`
	Operation     typeparse.QueryCondition `json:"operation"`
	OperationInfo typeparse.QueryCondition `json:"operation_info"`
	Quick         typeparse.QueryCondition `json:"quick"`
}
