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
	LocalDockerRepositoryPort = 8082
)

var (
	NotSupportedBatchOperation = errors.New("not supported operation")
)

var (
	ResourceDir = "resource"
	ChartsDir   = path.Join(ResourceDir, "charts")
)
