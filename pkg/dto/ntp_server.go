package dto

import "github.com/KubeOperator/KubeOperator/pkg/model"

type NtpServer struct {
	model.NtpServer
}

type NtpServerCreate struct {
	Name    string `json:"name" validate:"required"`
	Address string `json:"address" validate:"required"`
	Status  string `json:"status" validate:"required"`
}

type NtpServerUpdate struct {
	Address string `json:"address" validate:"required"`
	Status  string `json:"status" validate:"required"`
}
