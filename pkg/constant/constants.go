package constant

import (
	"errors"
	"path"
)

const (
	PageNumQueryKey  = "pageNum"
	PageSizeQueryKey = "pageSize"

	BatchOperationUpdate = "update"
	BatchOperationCreate = "create"
	BatchOperationDelete = "delete"

	LocalRpmRepositoryPort    = 8081
	LocalHelmRepositoryPort   = 8081
	LocalDockerRepositoryPort = 8082

	DefaultResourceName = "kubeoperator"
	StatusPending       = "pending"
	StatusRunning       = "running"
)

var (
	NotSupportedBatchOperation = errors.New("not supported operation")
)

var (
	ResourceDir          = "resource"
	ChartsDir            = path.Join(ResourceDir, "charts")
	DefaultDataDir       = "/var/ko/data"
	DefaultAnsibleLogDir = path.Join(DefaultDataDir, "ansible")
	BackupDir            = path.Join(DefaultDataDir, "backup")
	DefaultRepositoryDir = path.Join(DefaultDataDir, "git")
)
