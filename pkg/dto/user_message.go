package dto

import "github.com/KubeOperator/KubeOperator/pkg/model"

type UserMessageDTO struct {
	model.UserMessage
}

type UserMessageOp struct {
	Operation string           `json:"operation"`
	Items     []UserMessageDTO `json:"items"`
}
