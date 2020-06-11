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
)

var (
	NotSupportedBatchOperation = errors.New("not supported operation")
)

var (
	ResourceDir = "resource"
	ChartsDir   = path.Join(ResourceDir, "charts")
)
