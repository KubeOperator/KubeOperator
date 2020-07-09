package dto

import (
	"github.com/KubeOperator/KubeOperator/pkg/model"
	v1 "k8s.io/api/core/v1"
)

type Node struct {
	model.ClusterNode
	Info v1.Node `json:"info"`
}

type NodeCreation struct {
	Hosts []struct {
		Role     string `json:"role"`
		HostName string `json:"hostName"`
	} `json:"hosts"`
	Increase int `json:"increase"`
}
