package dto

import (
	"github.com/KubeOperator/KubeOperator/pkg/model"
)

type Host struct {
	model.Host
	ClusterName string `json:"clusterName"`
	ProjectName string `json:"projectName"`
	ZoneName    string `json:"zoneName"`
}

type HostCreate struct {
	Name         string `json:"name" validate:"required,max=30,hostname"`
	Ip           string `json:"ip" validate:"required,koip"`
	Port         int    `json:"port" validate:"required,gte=1,lte=65535"`
	CredentialID string `json:"credentialId" validate:"required"`
}

type HostPage struct {
	Items []Host `json:"items"`
	Total int    `json:"total"`
}

type HostOp struct {
	Operation string `json:"operation" validate:"required"`
	Items     []Host `json:"items" validate:"required"`
}

type HostSync struct {
	HostName   string `json:"hostName"`
	HostStatus string `json:"hostStatus"`
}
