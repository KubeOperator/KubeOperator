package dto

import (
	"github.com/KubeOperator/KubeOperator/pkg/model"
)

type TaskLog struct {
	model.TaskLog `json:"tasklogs"`
}

type Logs struct {
	Msg string `json:"msg"`
}
