package dto

import "github.com/KubeOperator/KubeOperator/pkg/model"

type SystemRegistry struct {
	model.SystemRegistry
	Status  string `json:"status"`
	Message string `json:"message"`
}

type SystemRegistryCreate struct {
	Hostname           string `json:"hostname" validate:"required"`
	Protocol           string `json:"protocol" validate:"required"`
	Architecture       string `json:"architecture" validate:"required"`
	RepoPort           int    `json:"repoPort" validate:"required"`
	RegistryPort       int    `json:"registryPort" validate:"required"`
	RegistryHostedPort int    `json:"registryHostedPort" validate:"required"`
	NexusUser          string `json:"nexusUser" validate:"required"`
	NexusPassword      string `json:"nexusPassword" validate:"required"`
}

type SystemRegistryUpdate struct {
	ID                 string `json:"id"`
	Hostname           string `json:"hostname"`
	Protocol           string `json:"protocol"`
	Architecture       string `json:"architecture"`
	RepoPort           int    `json:"repoPort"`
	RegistryPort       int    `json:"registryPort"`
	RegistryHostedPort int    `json:"registryHostedPort"`
}

type SystemRegistryDelete struct {
	Architecture string `json:"architecture" validate:"required"`
}

type SystemRegistryBatchOp struct {
	Operation string           `json:"operation" validate:"required"`
	Items     []SystemRegistry `json:"items" validate:"required"`
}

type RepoChangePassword struct {
	ID            string `json:"id"`
	NexusUser     string `json:"nexusUser"`
	NexusPassword string `json:"nexusPassword"`
}

type SystemRegistryConn struct {
	Hostname      string `json:"hostname" validate:"required"`
	Protocol      string `json:"protocol" validate:"required"`
	RepoPort      int    `json:"repoPort" validate:"required"`
	NexusUser     string `json:"nexusUser"`
	NexusPassword string `json:"nexusPassword"`
}
